// database.post
package database

import (
	"strconv"
	"strings"
)

// Struttura dati che memorizza sia l'utente la posizione del 
// suo cursore nell'attuale conversazione
type convUser struct {
	user    *User
	post    int
	cursore int
}

// Struttura dati che memorizza la conversazione
type Conversation struct {
	Title        string // Titolo della conversazione
	ID           int    // Identificativo numerico intero univoco della conversazione
	UtenteAttivo int    // ID dell'utente attivo sulla conversazione

	//TODO totale_numero_post int
	//TODO privata	         bool

	connected     map[int]*convUser // Mappa degli utenti attualmente connessi alla conversazione
	postMap       map[int]*Post     // Mappa dei post per un accesso più immediato ad essi
	contatorePost int               // Contatore di quanti post sono stati creati per poterne creare uno nuovo con id univoco
	testaPost     *Post             // Puntatore alla testa dell'albero dei post
}

// Errore nel gestire la conversazione
type ConversationError struct {
	errExplanation string
	conv           *Conversation
}

// Ritorna le posizioni di tutte le tuple di cursori sotto forma di stringa già pronta 
// essere spedita sulla rete
func (conv *Conversation) GetAllPositionString() string {
	var messaggio string
	posizioni := make([]string, 0, len(conv.connected))

	for k, usr := range conv.connected {
		messaggio = "\\U" + strconv.Itoa(k) + "\\\\P" + strconv.Itoa(usr.post) + "\\\\C" + strconv.Itoa(usr.cursore) + "\\"
		posizioni = append(posizioni, messaggio)
	}

	return strings.Join(posizioni, "")
}

// Crea un nuovo ConversationError
func NewConversationError(errExplanation string, conv *Conversation) *ConversationError {
	err := new(ConversationError)
	err.errExplanation = errExplanation
	err.conv = conv
	return err
}

// converte il ConversationError in una stringa esplicativa dell'errore
func (err *ConversationError) Error() string {
	return "Errore nella conversazione: " + err.conv.Title + "\n" + err.errExplanation
}

// Crea una nuova conversazione inizializzando tutti i dati 
func NewConversation(creator *User) *Conversation {
	conv := new(Conversation)

	//connected     map[int]*User
	//postMap       map[int]*Post
	conv.connected = make(map[int]*convUser)
	conv.postMap = make(map[int]*Post)

	conv.contatorePost = 0
	conv.connected[creator.ID] = &convUser{creator, 0, 0}
	conv.testaPost = conv.NewPost(creator, nil)
	//conv.testaPost.write(&ServerFakeUser,[]rune("Conversazione senza nome"),0)

	return conv
}

// Imposta lo stato online dell'utente "user". Se l'utente è già connesso 
// viene ritornato un errore 
func (conv *Conversation) NewUserConnection(user *User) *ConversationError {

	if userold, ok := conv.connected[user.ID]; ok {
		return NewConversationError("L'utente "+userold.user.Username+" è già connesso", conv)
	} else {
		conv.connected[user.ID] = &convUser{user, 0, 0}
		return nil
	}
	return nil
}

// Funzione per gestire la disconnessione del vari utenti
func (conv *Conversation) UserDisconnection(user *User) {
	delete(conv.connected, user.ID)
}

//create a post in response of the given post
func (conv *Conversation) CreateResponseToPost(idPost int, user *User) (*Post, error) {
	padre := conv.postMap[idPost]

	risposta, err := padre.Respond(conv, user)
	if err != nil {
		return nil, err
	}

	return risposta, nil
}

// Cambia la posizione del cursore all'interno del post selezionato
func (conv *Conversation) ChangePos(user *User, pos int) error {
	userconv, ok := conv.connected[user.ID]
	if ok == false {
		return NewConversationError("Utente "+user.Username+" che scrive non è connesso", conv)
	}
	//TODO controlloare che il nuovo cursore abbia senso
	userconv.cursore = pos
	return nil
}

// Cambia il post selezionato in cui scrivere
func (conv *Conversation) ChangePost(user *User, post int) error {
	userconv, ok := conv.connected[user.ID]
	if ok == false {
		return NewConversationError("Utente "+user.Username+" che scrive non è connesso", conv)
	}
	userconv.post = post
	return nil
}

// Aggiunge una parte di testo alla conversazione come azione di user
func (conv *Conversation) InsElem(user *User, appendRunes []rune) error {
	userconv, ok := conv.connected[user.ID]
	if ok == false {
		return NewConversationError("Utente "+user.Username+" che scrive non è connesso", conv)
	}

	post, ok := conv.postMap[userconv.post]
	if ok == false {
		return NewConversationError("Utente "+user.Username+" scrive in un post inesistente ("+strconv.Itoa(userconv.post)+")", conv)
	}

	post.write(user, appendRunes, userconv.cursore)
	userconv.cursore += len(appendRunes)
	return nil
}

// Elimina una parte di testo dalla conversazione come azione di user
func (conv *Conversation) DelElem(user *User, howmany int) error {
	userconv, ok := conv.connected[user.ID]
	if ok == false {
		return NewConversationError("Utente "+user.Username+" che cancella non è connesso", conv)
	}

	userconv.cursore -= howmany

	post, ok := conv.postMap[userconv.post]
	if ok == false {
		return NewConversationError("Utente "+user.Username+" cancella in un post inesistente ("+strconv.Itoa(post.idPost)+")", conv)
	}

	post.del(user, userconv.cursore, howmany)
	return nil
}

// Ritorna l'intera conversazione visualizzabile come un unica stringa
func (conv *Conversation) GetComplete(separators bool) string {
	//	return conv.testaPost.testo.GetComplete(separators)
	//TODO tornare una forma totale della stringa
	//	var totale []rune

	//	postAttuale := conv.testaPost
	//	for postAttuale != nil {
	//		totale = append(totale, []rune("\\P"+strconv.Itoa(postAttuale.idPost)+"\\")...)
	//		totale = append(totale, []rune(postAttuale.testo.GetComplete(separators))...)
	//		postAttuale = postAttuale.rispostaPrincipale
	//	}

	totale := conv.testaPost.getTestoPiuFigli(separators)

	return "\\P0\\\\C0\\" + string(totale)
}

//funzione ricorsiva per srotolare cosa c'è scritto dentro tutti i post in un'unica stringa
func (post *Post) getTestoPiuFigli(separators bool) []rune {
	var totale []rune = make([]rune, 0, 256)

	totale = append(totale, []rune(post.testo.GetComplete(separators))...)
	if post.rispostaPrincipale != nil {
		testoRispPrincipale := post.rispostaPrincipale.getTestoPiuFigli(separators)
		totale = append(totale, []rune("\\K"+strconv.Itoa(post.idPost)+"\\"+strconv.Itoa(post.rispostaPrincipale.idPost)+"\\")...)
		if !separators {
			totale = append(totale, []rune("\\P"+strconv.Itoa(post.rispostaPrincipale.idPost)+"\\"+"\\C0\\")...)
		}
		totale = append(totale, testoRispPrincipale...)
	}
	if post.rispostaSecondaria != nil {
		testoRispSecondaria := post.rispostaSecondaria.getTestoPiuFigli(separators)
		totale = append(totale, []rune("\\K"+strconv.Itoa(post.idPost)+"\\"+strconv.Itoa(post.rispostaSecondaria.idPost)+"\\")...)
		if !separators {
			totale = append(totale, []rune("\\P"+strconv.Itoa(post.rispostaSecondaria.idPost)+"\\"+"\\C0\\")...)
		}
		totale = append(totale, testoRispSecondaria...)
	}
	return totale
}

// Crea un post in risposta al post identificato da idParentPost, user viene immesso
// come scrittore iniziale del post
func (conv *Conversation) Respond(idParentPost int, user *User) (int, error) {
	_, ok := conv.connected[user.ID]
	if ok == false {
		return 0, NewConversationError("Utente "+user.Username+" che cancella non è connesso", conv)
	}

	post, ok := conv.postMap[idParentPost]
	if ok == false {
		return 0, NewConversationError("Utente "+user.Username+" cancella in un post inesistente ("+strconv.Itoa(post.idPost)+")", conv)
	}

	postR, err := post.Respond(conv, user)
	if err != nil {
		return 0, err
	}

	return postR.idPost, nil
}

/*
func (post *Post) Respond(conv *Conversation, user *User) (*Post, error) {
	response := conv.NewPost(user, post)
	if post.rispostaPrincipale == nil {
		post.rispostaPrincipale = response
		return response, nil
	}
	if post.rispostaSecondaria == nil {
		post.rispostaSecondaria = response
		return response, nil
	}
	return nil, NewConversationError("Impossibile attaccare più di due risposte ad un solo post", conv)
}
*/
