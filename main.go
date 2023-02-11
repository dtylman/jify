package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"os/exec"
	"strconv"

	"github.com/dtylman/console"
	"github.com/dustin/go-humanize"
	"github.com/iancoleman/strcase"
	"github.com/kkdai/youtube/v2"
)

type Jif struct {
	VideoID  string   `json:"video_id"`
	Source   string   `json:"source"`
	Name     string   `json:"name"`
	File     string   `json:"file"`
	Start    float64  `json:"start"`
	Duration float64  `json:"duration"`
	Labels   []string `json:"labels"`
}

func (j *Jif) SetLabel(label string) {
	if label == "" {
		return
	}
	if j.Labels == nil {
		j.Labels = make([]string, 0)
	}
	j.Labels = append(j.Labels, label)
}

func work() error {

	videoID, err := console.Ask("video id?")
	if err != nil {
		return err
	}

	scale, err := askInt("scale", 320)
	if err != nil {
		return err
	}

	client := youtube.Client{}
	log.Printf("getting video %v", videoID)
	video, err := client.GetVideo(videoID)
	if err != nil {
		return err
	}

	type videoFormat struct {
		Index         int
		Quality       string
		QualityLabel  string
		FPS           int
		Width         int
		Height        int
		AudioChannels int
	}

	videoFormats := make([]videoFormat, 0)
	bestDelta := math.MaxInt
	best := 0
	for i, vf := range video.Formats {
		videoFormats = append(videoFormats,
			videoFormat{Index: i,
				Quality:       vf.Quality,
				QualityLabel:  vf.QualityLabel,
				FPS:           vf.FPS,
				Width:         vf.Width,
				Height:        vf.Height,
				AudioChannels: vf.AudioChannels,
			})
		delta := math.Abs(float64(vf.Width - scale))
		if bestDelta > int(delta) {
			bestDelta = int(delta)
			best = i
		}
	}
	console.WriteTable(videoFormats)
	formatID, err := askInt("id?", best)
	if err != nil {
		return err
	}

	startingPoint, err := askFloat("staring time", 0)
	if err != nil {
		return err
	}

	duration := video.Duration.Seconds() - startingPoint
	duration, err = askFloat("duration", duration)
	if err != nil {
		return err
	}

	speed, err := askInt("speed", 100)
	if err != nil {
		return err
	}
	pts := 100 / speed

	fps, err := askInt("fps", 5)
	if err != nil {
		return err
	}

	log.Println("getting video stream")
	stream, size, err := client.GetStream(video, &video.Formats[formatID])
	if err != nil {
		return err
	}
	defer stream.Close()

	inputFileName := strcase.ToSnake(videoID) + ".vid"

	_, err = os.Stat(inputFileName)
	if err != nil {
		file, err := os.Create(inputFileName)
		if err != nil {
			return err
		}

		defer func() {
			file.Close()
		}()

		log.Printf("downloading %v bytes to %v", humanize.Bytes(uint64(size)), file.Name())
		_, err = io.Copy(file, stream)
		if err != nil {
			return err
		}
	}

	name, err := console.AskOptions("name", transliterate(strcase.ToSnake(video.Title)), true, false)
	if err != nil {
		return err
	}

	outFileName := transliterate(name) + ".gif"
	videoFilter := fmt.Sprintf("fps=%v,setpts=%v*PTS,scale=%v:-1:flags=lanczos,split[s0][s1];[s0]palettegen[p];[s1][p]paletteuse",
		fps, pts, scale)
	cmd := exec.Command(
		"ffmpeg",
		"-i", inputFileName,
		"-ss", fmt.Sprintf("%v", startingPoint),
		"-t", fmt.Sprintf("%v", duration),
		"-vf", videoFilter,
		"-loop", "0",
		outFileName)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	log.Println(cmd)
	err = cmd.Run()
	if err != nil {
		return err
	}

	jif := Jif{
		VideoID:  videoID,
		Source:   fmt.Sprintf("https://www.youtube.com/watch?v=%v", videoID),
		Name:     name,
		File:     outFileName,
		Start:    startingPoint,
		Duration: duration,
	}
	jif.SetLabel(name)
	jif.SetLabel(video.Title)
	jif.SetLabel(video.Author)
	jif.SetLabel(video.Description)

	data, err := json.MarshalIndent(jif, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(name+".json", data, 0755)
}

func askInt(prompt string, def int) (int, error) {
	answer, err := console.AskOptions(prompt, strconv.Itoa(def), true, false)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(answer)

}

func askFloat(prompt string, def float64) (float64, error) {
	answer, err := console.AskOptions(prompt, fmt.Sprintf("%v", def), true, false)
	if err != nil {
		return 0, err
	}
	return strconv.ParseFloat(answer, 64)

}
func main() {
	err := work()
	if err != nil {
		fmt.Println(err)
	}
}
