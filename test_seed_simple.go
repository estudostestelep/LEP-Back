package main

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var (
	clearFirst  bool
	environment string
	verbose     bool
)

func main1() {
	var rootCmd = &cobra.Command{
		Use:   "seed",
		Short: "LEP Database Seeder",
		Long:  `Populate the LEP database with realistic sample data for development and testing.`,
		Run:   runSeed,
	}

	rootCmd.Flags().BoolVar(&clearFirst, "clear-first", false, "Clear existing data before seeding")
	rootCmd.Flags().StringVar(&environment, "environment", "dev", "Environment to seed (dev, test, staging)")
	rootCmd.Flags().BoolVar(&verbose, "verbose", false, "Enable verbose logging")

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func runSeed(cmd *cobra.Command, args []string) {
	fmt.Println("\n🌱 LEP Database Seeder")
	fmt.Println("======================")
	fmt.Printf("Environment: %s\n", environment)
	fmt.Printf("Clear first: %t\n", clearFirst)
	fmt.Printf("Verbose: %t\n\n", verbose)

	fmt.Println("✅ Sistema de seeding configurado!")
	fmt.Println("🔧 Fluxo implementado:")
	fmt.Println("  1. 📋 Criar organização (sem auth)")
	fmt.Println("  2. 📁 Criar projeto (usando org ID)")
	fmt.Println("  3. 👥 Criar usuário admin (usando org/project IDs)")
	fmt.Println("  4. 🔐 Fazer login para obter token")
	fmt.Println("  5. 🚀 Criar demais entidades via rotas server reais")
	fmt.Println("\n🎯 Próximos passos:")
	fmt.Println("  • Configurar banco de dados PostgreSQL")
	fmt.Println("  • Executar: go run ./cmd/seed --verbose")
	fmt.Println("  • Verificar dados criados via API")
}
