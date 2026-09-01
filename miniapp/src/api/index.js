import request from '../utils/request'

export function register(data) {
  return request({ url: '/register', method: 'POST', data })
}

export function login(data) {
  return request({ url: '/login', method: 'POST', data })
}

export function getUserInfo() {
  return request({ url: '/user/info' })
}

export function queryWord(data) {
  return request({ url: '/words/query', method: 'POST', data })
}

export function saveWord(data) {
  return request({ url: '/words/save', method: 'POST', data })
}

export function getWords(params) {
  return request({ url: '/words', data: params })
}

export function deleteWord(id) {
  return request({ url: '/words/' + id, method: 'DELETE' })
}
