package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func init() {
	programName := "create-repo"
	flag.Usage = func() {
		fmt.Printf("\nUsage: %s [options]\n", programName)
		fmt.Printf("\nOptions:\n")
		fmt.Printf("\t-name                 Repository name (required). If remote repository already exists, ask whether to clone it\n")
		fmt.Printf("\t-create-dir           Create a directory for the repository (default: current directory). If the directory already exists, ask whether to reuse it\n")
		fmt.Printf("\t-private              Make the repository private (default: public).\n")
		fmt.Printf("\t-help                 Show this help message\n")
		fmt.Printf("\nExamples:\n")
		fmt.Printf("\t%s -help\n", programName)
		fmt.Printf("\t%s -name=myrepo\n", programName)
		fmt.Printf("\t%s -name=myrepo -create-dir=true\n", programName)
		fmt.Printf("\t%s -name=myrepo -private=true\n", programName)
		fmt.Printf("\t%s -name=myrepo -create-dir=true -private=true\n\n", programName)

		fmt.Println("\n🔥 Automates repository creation with Git and GitHub CLI 🔥")
	}
}

func main() {
	pRepoName := flag.String("name", "", "Repository name (required)")
	createDir := flag.Bool("create-dir", false, "Create a directory for the repository (default is current directory)")
	private := flag.Bool("private", false, "Make the repository private (default is public)")
	help := flag.Bool("help", false, "Show this help message")

	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	if *pRepoName == "" {
		fmt.Println("❌ Error: Repository name is required.")
		flag.Usage()
		os.Exit(1)
	}

	repoName := *pRepoName

	fmt.Printf("⌛ Checking if repository '%s' already exists...\n", repoName)

	// 🚀 Verifica se o repositório remoto já existe
	if repoExists(repoName) {
		// 🔥 Pergunta se deseja clonar em vez de criar um novo
		if confirm(fmt.Sprintf("⚠️  Repository '%s' already exists on GitHub.\n ▶️ Do you want to clone it instead?", repoName)) {
			runCommand("git", "clone", "https://github.com/YOUR_USERNAME/"+repoName+".git")
			fmt.Println("✅ Repository cloned successfully!")
			os.Exit(0) // Encerra o programa após o clone
		}

		// Se não quiser clonar, perguntar se quer sobrescrever
		if !confirm(fmt.Sprintf("⚠️  Warning: Repository '%s' already exists.\n ▶️ Do you want to continue and use the existing repository?", repoName)) {
			fmt.Println("❌ Operation canceled.")
			os.Exit(1)
		}

		fmt.Printf("📡 Using existing remote repository: %s\n", repoName)
	} else {
		fmt.Printf("🚀 Repository '%s' does not exist on GitHub. Proceeding with creation...\n", repoName)
	}

	if *createDir {
		createDirectory(repoName)
	}

	runCommand("git", "init")
	runCommand("git", "add", ".")
	runCommand("git", "commit", "--allow-empty", "-m", "ci: create repository")

	visibility := "--public"
	if *private {
		visibility = "--private"
	}

	runCommand("gh", "repo", "create", repoName, visibility, "--source=.", "--remote=origin", "--push")

	fmt.Printf("✅ Repository '%s' created successfully and synchronized with Github!\n", repoName)
}

// 🚀 Verifica se o repositório remoto já existe
func repoExists(repoName string) bool {
	cmd := exec.Command("gh", "repo", "view", repoName)
	err := cmd.Run()
	return err == nil // Se não houver erro, significa que o repositório já existe
}

// 🚀 Cria o diretório localmente se necessário
func createDirectory(repoName string) {
	if _, err := os.Stat(repoName); err == nil {
		if !confirm(fmt.Sprintf("⚠️  Warning: Directory 📂 '%s' already exists.\n ▶️ Do you want to use it?", repoName)) {
			fmt.Println("❌ Operation canceled.")
			os.Exit(1)
		}
		fmt.Printf("📂 Using existing repository directory: %s\n", repoName)
	} else {
		if err := os.Mkdir(repoName, 0755); err != nil {
			log.Fatalf("❌ Error creating directory: %v", err)
		}
		fmt.Printf("✅ Repository directory '%s' created successfully!\n", repoName)
	}

	if err := os.Chdir(repoName); err != nil {
		log.Fatalf("❌ Error changing directory: %v", err)
	}
}

// 🚀 Função de confirmação para perguntar ao usuário
func confirm(prompt string) bool {
	fmt.Print(prompt + " (y/N): ")
	var response string
	fmt.Scanln(&response)
	return strings.ToLower(response) == "y"
}

// 🚀 Executa um comando do sistema
func runCommand(command string, args ...string) {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Printf("\t⚙️  Running command: %s %s\n", command, strings.Join(args, " "))

	err := cmd.Run()
	if err != nil {
		log.Fatalf("❌ Error running command %s %v: %v\n", command, args, err)
		os.Exit(1)
	}
}
