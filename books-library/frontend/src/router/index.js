import { createRouter, createWebHistory } from 'vue-router'
import Body from './../components/Body.vue'
import Login from './../components/Login.vue'
import Books from './../components/Books.vue'
import Book from './../components/Body.vue'
import BookEdit from './../components/BookEdit.vue'
import BookAdmin from './../components/BooksAdmin.vue'

const routes = [
  {
    path: '/',
    name: 'Home',
    component: Body,
  },
  {
    path: '/login',
    name: 'Login',
    component: Login,
  },
  {
    path: '/books',
    name: 'Books',
    component: Books,
  },
  {
    path: `/books/:bookName`,
    name: 'Book',
    component: Book,
  },
  {
    path: '/admin/books',
    name: 'BookAdmin',
    component: BookAdmin,
  },
  {
    path: '/admin/books/:bookId',
    name: 'BookEdit',
    component: BookEdit,
  },
]

const router = createRouter({ history: createWebHistory(), routes })
export default router
