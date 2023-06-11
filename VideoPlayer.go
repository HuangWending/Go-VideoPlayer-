package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/nsf/termbox-go"
)

const (
	playbackSpeed = 1.0 // 默认播放速度
)

var (
	videoFile      string // 视频文件路径
	isPlaying      bool   // 是否正在播放
	isPaused       bool   // 是否已暂停
	currentSpeed   float64
	volumeLevel    int // 音量级别
	audioSupported bool
)

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	// 初始化播放器状态
	isPlaying = false
	isPaused = false
	currentSpeed = playbackSpeed
	volumeLevel = 50

	// 获取视频文件路径
	videoFile = getVideoFilePath()

	// 检查是否支持音频
	audioSupported = checkAudioSupport()

	// 显示播放器界面
	render()

	// 监听键盘事件
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			if ev.Key == termbox.KeyEsc {
				return
			}

			switch ev.Ch {
			case 'p', 'P': // 播放/暂停
				togglePlayPause()
			case 's', 'S': // 减速
				decreaseSpeed()
			case 'f', 'F': // 加速
				increaseSpeed()
			case 'u', 'U': // 音量增加
				increaseVolume()
			case 'd', 'D': // 音量减少
				decreaseVolume()
			}

			render()
		case termbox.EventError:
			panic(ev.Err)
		}
	}
}

// 获取视频文件路径
func getVideoFilePath() string {
	fmt.Print("请输入视频文件的路径：")
	var filePath string
	fmt.Scanln(&filePath)

	return filePath
}

// 检查是否支持音频
func checkAudioSupport() bool {
	cmd := exec.Command("ffprobe", "-show_streams", "-i", videoFile)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("检查音频支持时出错:", err)
		return false
	}

	streams := strings.Split(string(output), "[STREAM]")
	for _, stream := range streams {
		if strings.Contains(stream, "codec_type=audio") {
			return true
		}
	}

	return false
}

// 播放/暂停切换
func togglePlayPause() {
	if !isPlaying && !isPaused {
		playVideo()
	} else if isPlaying && !isPaused {
		pauseVideo()
	} else if isPlaying && isPaused {
		resumeVideo()
	}
}

// 播放视频
func playVideo() {
	cmd := exec.Command("ffplay", "-i", videoFile, "-vf", fmt.Sprintf("setpts=%.2f*PTS", currentSpeed))
	err := cmd.Start()
	if err != nil {
		fmt.Println("播放视频时出错:", err)
		return
	}

	isPlaying = true
	isPaused = false
}

// 暂停视频
func pauseVideo() {
	cmd := exec.Command("pkill", "-STOP", "-f", "ffplay")
	err := cmd.Run()
	if err != nil {
		fmt.Println("暂停视频时出错:", err)
		return
	}

	isPaused = true
}

// 恢复播放视频
func resumeVideo() {
	cmd := exec.Command("pkill", "-CONT", "-f", "ffplay")
	err := cmd.Run()
	if err != nil {
		fmt.Println("恢复视频播放时出错:", err)
		return
	}

	isPaused = false
}

// 增加播放速度
func increaseSpeed() {
	currentSpeed += 0.1
	if currentSpeed > 2.0 {
		currentSpeed = 2.0
	}
}

// 减少播放速度
func decreaseSpeed() {
	currentSpeed -= 0.1
	if currentSpeed < 0.1 {
		currentSpeed = 0.1
	}
}

// 增加音量
func increaseVolume() {
	volumeLevel += 10
	if volumeLevel > 100 {
		volumeLevel = 100
	}
}

// 减少音量
func decreaseVolume() {
	volumeLevel -= 10
	if volumeLevel < 0 {
		volumeLevel = 0
	}
}

// 绘制界面
func render() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	// 绘制视频播放器框架
	drawPlayerFrame()

	// 绘制功能按钮
	drawControls()

	termbox.Flush()
}

// 绘制视频播放器框架
func drawPlayerFrame() {
	width, height := termbox.Size()

	for y := 0; y < height; y++ {
		termbox.SetCell(0, y, '│', termbox.ColorDefault, termbox.ColorDefault)
		termbox.SetCell(width-1, y, '│', termbox.ColorDefault, termbox.ColorDefault)
	}

	for x := 1; x < width-1; x++ {
		termbox.SetCell(x, 0, '─', termbox.ColorDefault, termbox.ColorDefault)
		termbox.SetCell(x, height-1, '─', termbox.ColorDefault, termbox.ColorDefault)
	}

	termbox.SetCell(0, 0, '┌', termbox.ColorDefault, termbox.ColorDefault)
	termbox.SetCell(width-1, 0, '┐', termbox.ColorDefault, termbox.ColorDefault)
	termbox.SetCell(0, height-1, '└', termbox.ColorDefault, termbox.ColorDefault)
	termbox.SetCell(width-1, height-1, '┘', termbox.ColorDefault, termbox.ColorDefault)
}

// 绘制功能按钮
func drawControls() {
	width, _ := termbox.Size()
	y := 2

	drawButton(2, y, "[P] 播放/暂停", !isPlaying || isPaused)
	y += 2

	drawButton(2, y, "[S] 减速", isPlaying && !isPaused && currentSpeed <= 0.1)
	drawButton(16, y, "[F] 加速", isPlaying && !isPaused && currentSpeed >= 2.0)
	y += 2

	drawButton(2, y, "[U] 音量增加", !isPlaying || isPaused || volumeLevel >= 100)
	drawButton(18, y, "[D] 音量减少", !isPlaying || isPaused || volumeLevel <= 0)
	y += 2

	if !audioSupported {
		drawMessage(2, y, "警告: 视频不包含音频轨道")
	}
}

// 绘制按钮
func drawButton(x, y int, label string, disabled bool) {
	fgColor := termbox.ColorDefault
	bgColor := termbox.ColorGreen

	if disabled {
		fgColor = termbox.ColorBlack
		bgColor = termbox.ColorRed
	}

	for _, ch := range label {
		termbox.SetCell(x, y, ch, fgColor, bgColor)
		x++
	}
}

// 绘制消息
func drawMessage(x, y int, message string) {
	fgColor := termbox.ColorRed
	bgColor := termbox.ColorDefault

	for _, ch := range message {
		termbox.SetCell(x, y, ch, fgColor, bgColor)
		x++
	}
}
