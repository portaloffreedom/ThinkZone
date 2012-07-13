// database.post
package database

import (
	"container/list"

//	"crypto/sha256"

//	"fmt"
)

type ConversationError struct {
	errExplanation string
	conv           *Conversation
}

func NewConversationError(errExplanation string, conv *Conversation) *ConversationError {
	err := new(ConversationError)
	err.errExplanation = errExplanation
	err.conv = conv
	return err
}

func (err *ConversationError) Error() string {
	return "Errore nella conversazione: " + err.conv.Title + "\n" + err.errExplanation
}

type Conversation struct {
	Title            string
	ID_conversazione int

	//totale_numero_post int
	//privata	         bool

	connected     map[int]*User
	postMap       map[int]*Post
	contatorePost int
	TestaPost     *Post
}

type Post struct {
	testo   *SuperString
	idPost  int
	writers *list.List //list of users

	padre    *Post      //puntatore a cosa sto rispondendo
	risposte *list.List //puntatore alle risposte
}

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

func (post *Post) Text(user *User) *SuperString {
	post.addWriter(user)
	return post.testo
}

func (conv *Conversation) TotalConversation() string {
	return conv.TestaPost.testo.GetComplete(false)
}

func (post *Post) Respond(conv *Conversation, user *User) *Post {
	response := conv.NewPost(user, post)
	post.risposte.PushBack(response)
	return response
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

func NewConversation(creator *User) *Conversation {
	conv := new(Conversation)

	//connected     map[int]*User
	//postMap       map[int]*Post
	conv.connected = make(map[int]*User)
	conv.postMap = make(map[int]*Post)

	conv.contatorePost = 1
	conv.connected[creator.ID] = creator
	conv.TestaPost = conv.NewPost(creator, nil)

	return conv
}

func (conv *Conversation) NewUserConnection(user *User) *ConversationError {

	if userold, ok := conv.connected[user.ID]; ok {
		return NewConversationError("L'utente "+userold.Username+" è già connesso", conv)
	} else {
		conv.connected[user.ID] = user
		return nil
	}
	return nil
}

func (conv *Conversation) UserDisconnection(user *User) {
	delete(conv.connected, user.ID)
}

//create a post in response of the given post
func (conv *Conversation) ResponseToPost(idPost int, user *User) *Post {
	padre := conv.postMap[idPost]
	risposta := padre.Respond(conv, user)

	return risposta
}
