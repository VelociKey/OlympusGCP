package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	fmt.Println("Starting OlympusGCP-FinOps Provisioner (Dagger-Driven)...")

	// Detect Fleet Root
	wd, _ := os.Getwd()
	root := wd
	if filepath.Base(wd) == "OlympusGCP-FinOps" {
		root = filepath.Dir(wd)
	}

	forgePkg := filepath.Join(root, "00SDLC", "OlympusForge", "90000-Enablement-Labs", "900-Forge")

	fmt.Println("🔨 Building OlympusGCP-FinOps via Forge Pipeline...")

	cmd := exec.Command("go", "run", forgePkg, "-target", "native", "-workspace", "00SDLC/OlympusGCP-FinOps")
	cmd.Dir = root
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("❌ OlympusGCP-FinOps provisioning failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✅ OlympusGCP-FinOps provisioning complete.")
}
