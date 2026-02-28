package main

import (
	"fmt"
	"github.com/kkdai/youtube/v2"
	"github.com/mikkyang/id3-go"
	"io"
	"os"
	"os/exec"
)

type getErr struct {
	msg string
}

func (err *getErr) Error() string {
	return err.msg
}

func newGetErr() *getErr {
	return &getErr{
		msg: "couldn't download video from the link. try again",
	}
}

func DownloadAudioFromVideo(log Logger, link string) (string, error) {
	dwnldErr := newGetErr()

	client := youtube.Client{}
	video, err := client.GetVideo(link)
	if err != nil {
		log.Error(fmt.Sprintf("%v", dwnldErr))
		return "", dwnldErr
	}

	formats := video.Formats.WithAudioChannels()
	stream, _, err := client.GetStream(video, &formats[0])
	if err != nil {
		log.Error(fmt.Sprintf("error during client.GetStream: %v", err))
		return "", dwnldErr
	}

	defer errClose(stream.Close, log)

	vidTitle := video.Title

	tempVid := createName(vidTitle, ".mp4")
	err = createFile(stream, tempVid, log)

	defer func() {
		delErr := deleteFile(tempVid, log)
		if delErr != nil {
			log.Error(fmt.Sprintf("error during temp video file deletion: %v", err))
		}
	}()

	audioFileName, err := convertVideoToAudio(tempVid, vidTitle, log)
	if err != nil {
		log.Error(fmt.Sprintf("error during file format conversion: %v", err))
		return "", dwnldErr
	}

	return audioFileName, nil

}

func convertVideoToAudio(videoPath, videoTitle string, log Logger) (string, error) {
	audioFileName := createName(videoTitle, ".mp3")
	err := convertToMP3(videoPath, audioFileName)
	if err != nil {
		return "", err
	}

	file, err := id3.Open(audioFileName)
	if err != nil {
		return "", fmt.Errorf("error during metadata editing: %v", err)
	}

	defer errClose(file.Close, log)

	file.SetTitle(videoTitle)
	file.SetArtist("")

	return audioFileName, nil
}

func createFile(stream io.ReadCloser, filename string, log Logger) error {
	dwnldErr := newGetErr()

	file, err := os.Create(filename)
	if err != nil {
		log.Error(fmt.Sprintf("error during os.Create: %v", err))
		return dwnldErr
	}

	defer errClose(file.Close, log)

	_, err = io.Copy(file, stream)
	if err != nil {
		log.Error(fmt.Sprintf("error during io.Copy: %v", err))
		return dwnldErr
	}

	return nil
}

func deleteFile(filename string, log Logger) error {
	err := os.Remove(filename)
	if err != nil {
		log.Error(fmt.Sprintf("error during os.Create: %v", err))
		return err
	}

	return nil
}

func errClose(closerFunc func() error, log Logger) {
	err := closerFunc()
	if err != nil {
		log.Error(fmt.Sprintf("%v", err))
	}
}

func createName(name, format string) string {
	return fmt.Sprintf("%s.%s", name, format)
}

func convertToMP3(input string, output string) error {
	cmd := exec.Command("ffmpeg", "-i", input, "-vn", "-acodec", "libmp3lame", output)
	err := cmd.Run()
	return err
}
