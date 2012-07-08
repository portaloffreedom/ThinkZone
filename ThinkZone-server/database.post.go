// database.post
package main

import (
	//	"fmt"
	"container/list"
)

type ConversationError struct {
	errExplanation string
	conv           *Conversation
}

func (err *ConversationError) Error() string {
	return "Errore nella conversazione: " + err.conv.title + "\n" + err.errExplanation
}

type Conversation struct {
	title            string
	id_conversazione int

	//totale_numero_post int
	//privata	         bool

	connected     map[int]*User
	postMap       map[int]*Post
	contatorePost int
	testaPost     *Post
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

func (post *Post) Respond(conv *Conversation, user *User) *Post {
	response = conv.NewPost(user, post)
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

	conv.contatorePost = 1
	conv.connected[creator.id] = creator
	conv.testaPost = conv.NewPost(creator, nil, nil, nil)
}

func (conv *Conversation) NewUserConnection(user *User) *ConversationError {

	if conv.connected[user.id] != nil {
		return ConversationError{"L'utente " + user.username + " è già connesso", conv}
	}

	conv.connected[user.id] = user
	return nil

}

//create a post in response of the given post
func (conv *Conversation) ResponseToPost(idPost int, user *User) *Post {
	//TODO Conversation.ResponseToPost

	padre := conv.postMap[idPost]
	risposta = padre.Respond(conv, user)

	return risposta
}
