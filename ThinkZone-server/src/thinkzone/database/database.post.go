package database

import (
	"container/list"
)

type Post struct {
	testo   *SuperString
	idPost  int
	writers *list.List //list of users

	padre              *Post //puntatore a cosa sto rispondendo
	rispostaPrincipale *Post
	rispostaSecondaria *Post
	//	risposte *list.List //puntatore alle risposte
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

//deprecated
func (post *Post) text(user *User) *SuperString {
	post.addWriter(user)
	return post.testo
}

func (post *Post) write(user *User, appendRunes []rune, pos int) {
	post.text(user).InsElem(appendRunes, pos)
}
func (post *Post) del(user *User, pos int, howmany int) {
	post.text(user).DelElem(pos, howmany)
}

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
