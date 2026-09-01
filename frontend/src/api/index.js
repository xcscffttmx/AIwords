import request from './request'

export function register(data) {
  return request.post('/register', data)
}

export function login(data) {
  return request.post('/login', data)
}

export function getUserInfo() {
  return request.get('/user/info')
}

export function queryWord(data) {
  return request.post('/words/query', data)
}

export function saveWord(data) {
  return request.post('/words/save', data)
}

export function getWords(params) {
  return request.get('/words', { params })
}

export function deleteWord(id) {
  return request.delete(`/words/${id}`)
}
