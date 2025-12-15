func playAudioFile(filename string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin": // macOS
		cmd = exec.Command("afplay", filename)
	case "linux":
		// Try different players in order of preference
		players := []string{"mpg123", "ffplay", "mplayer", "play"}
		var playerFound bool
		for _, player := range players {
			if _, err := exec.LookPath(player); err == nil {
				switch player {
				case "ffplay":
					cmd = exec.Command(player, "-nodisp", "-autoexit", "-t", "10", filename)
				case "mpg123":
					cmd = exec.Command(player, "-q", filename)
				default:
					cmd = exec.Command(player, filename)
				}
				playerFound = true
				break
			}
		}
		if !playerFound {
			return fmt.Errorf("no audio player found. Please install mpg123, ffplay, mplayer, or sox")
		}
	case "windows":
		// Windows: try to use a command-line player
		// Note: PowerShell's SoundPlayer doesn't work well with MP3s
		// Users should install mpg123 or ffplay
		if _, err := exec.LookPath("ffplay"); err == nil {
			cmd = exec.Command("ffplay", "-nodisp", "-autoexit", "-t", "10", filename)
		} else if _, err := exec.LookPath("mpg123"); err == nil {
			cmd = exec.Command("mpg123", "-q", filename)
		} else {
			return fmt.Errorf("no audio player found. Please install ffplay or mpg123 for Windows")
		}
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	// Run the command with timeout - Fixed: proper error handling
	done := make(chan error, 1)

	go func() {
		// Start the command and wait for completion
		err := cmd.Run()
		done <- err
	}()

	select {
	case err := <-done:
		if err != nil {
			return fmt.Errorf("audio playback failed: %v", err)
		}
		return nil
	case <-time.After(12 * time.Second):
		// Timeout after 12 seconds
		if cmd.Process != nil {
			cmd.Process.Kill()
		}
		return nil // Return nil since timeout is expected for long files
	}
}
