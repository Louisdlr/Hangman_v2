package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

// Charger les mots depuis un fichier texte
func chargerMots(fichier string) []string {
	file, err := os.Open(fichier)
	if err != nil {
		fmt.Println("Erreur lors de l'ouverture du fichier :", err)
		os.Exit(1)
	}
	defer file.Close()

	var mots []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		mots = append(mots, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Erreur lors de la lecture du fichier :", err)
		os.Exit(1)
	}

	return mots
}

// Charger les étapes du pendu (images) depuis un fichier texte
func chargerPendu(fichier string) []string {
	file, err := os.Open(fichier)
	if err != nil {
		fmt.Println("Erreur lors de l'ouverture du fichier :", err)
		os.Exit(1)
	}
	defer file.Close()

	var etapes []string
	var etape string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		ligne := scanner.Text()
		if ligne == "" {
			etapes = append(etapes, etape)
			etape = ""
		} else {
			etape += ligne + "\n"
		}
	}
	if etape != "" {
		etapes = append(etapes, etape)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Erreur lors de la lecture du fichier :", err)
		os.Exit(1)
	}

	return etapes
}

// Choisir un mot aléatoirement dans la liste
func choisirMot(mots []string) string {
	rand.Seed(time.Now().UnixNano())
	return mots[rand.Intn(len(mots))]
}

// Révéler certaines lettres du mot dès le début
func revelerLettres(mot string, n int) []rune {
	motRevele := make([]rune, len(mot))
	for i := range motRevele {
		motRevele[i] = '_'
	}

	indicesReveles := rand.Perm(len(mot))[:n]
	for _, indice := range indicesReveles {
		motRevele[indice] = rune(mot[indice])
	}
	return motRevele
}

// Afficher le mot avec les lettres révélées et des espaces pour celles cachées
func afficherMotRevele(motRevele []rune) string {
	return strings.Join(strings.Split(string(motRevele), ""), " ")
}

// Afficher l'état du pendu à chaque étape
func afficherPendu(etapes []string, nbEssais int) {
	if nbEssais < 0 || nbEssais > 10 {
		nbEssais = 0
	}
	fmt.Println(etapes[10-nbEssais])
}

// Gérer la logique du jeu du pendu
func jouerPendu(mot string, etapes []string) {
	nbEssais := 10
	motRevele := revelerLettres(mot, len(mot)/2-1)
	lettresEssayees := make(map[rune]bool)

	for nbEssais > 0 {
		// Afficher le pendu et l'état actuel du mot
		afficherPendu(etapes, nbEssais)
		fmt.Println("Mot à deviner :", afficherMotRevele(motRevele))
		fmt.Println("Essais restants :", nbEssais)
		fmt.Print("Entrez une lettre : ")
		var lettre string
		fmt.Scanln(&lettre)
		lettre = strings.ToLower(lettre)
		if len(lettre) != 1 || !strings.Contains("abcdefghijklmnopqrstuvwxyz", lettre) {
			fmt.Println("Veuillez entrer une seule lettre valide.")
			continue
		}
		char := rune(lettre[0])
		if lettresEssayees[char] {
			fmt.Println("Vous avez déjà essayé cette lettre.")
			continue
		}
		lettresEssayees[char] = true

		// Si la lettre est dans le mot
		if strings.Contains(mot, lettre) {
			fmt.Println("Bravo ! La lettre", lettre, "est dans le mot.")
			for i, lettreMot := range mot {
				if lettreMot == char {
					motRevele[i] = char
				}
			}
		} else {
			fmt.Println("Dommage, la lettre", lettre, "n'est pas dans le mot.")
			nbEssais--
		}

		// Si le joueur a trouvé le mot
		if string(motRevele) == mot {
			fmt.Println("Félicitations ! Vous avez trouvé le mot :", mot)
			return
		}
	}

	// Si le joueur a perdu
	afficherPendu(etapes, 0)
	fmt.Println("Désolé, vous avez perdu. Le mot était :", mot)
}

// Fonction principale
func main() {
	mots := chargerMots("words.txt")
	etapes := chargerPendu("hangman.txt")
	mot := choisirMot(mots)
	jouerPendu(mot, etapes)
}
