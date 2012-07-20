package database

import (
	"container/list"
)

// Struttura dati per la memorizzazione del singolo post
type Post struct {
	testo  *SuperString //testo contenuto nel post
	idPost int          //Identificativo numerico intero univoco del post 

	//lista di tutti gli utenti che hanno partecipato alla creazione di questo post
	writers *list.List

	padre              *Post //puntatore a cosa sto rispondendo
	rispostaPrincipale *Post //puntatore al post risposta principale
	rispostaSecondaria *Post //puntatore al post risposta secondaria
}

//Crea una nuova struttura dati Post con i valori inizializzati 
func (conv *Conversation) NewPost(creator *User, padre *Post) *Post {
	post := new(Post)

	post.testo = NewSuperString()
	post.idPost = conv.contatorePost
	post.writers = list.New()
	post.writers.PushFront(creator)
	conv.contatorePost++
	conv.postMap[post.idPost] = post

	return post
}

// accede al testo del post memorizzando se è stato aggiunto
// un altro scrittore al post 
func (post *Post) text(user *User) *SuperString {
	post.addWriter(user)
	return post.testo
}

// aggiunge dei caratteri al testo del post
func (post *Post) write(user *User, appendRunes []rune, pos int) {
	post.text(user).InsElem(appendRunes, pos)
}

// elimina dei caratteri al testo del post
func (post *Post) del(user *User, pos int, howmany int) {
	post.text(user).DelElem(pos, howmany)
}

// crea un post di risposta al post invocato (con la struttura dati già
// inizializzata per bene)
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

//return false if the User was already a Writer of this post
func (post *Post) addWriter(writer *User) bool {

	for e := post.writers.Front(); e != nil; e = e.Next() {
		if e.Value == writer {
			return false
		}
	}

	post.writers.PushBack(writer)
	return true
}
