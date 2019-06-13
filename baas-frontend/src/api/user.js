import request from '@/utils/request'

export function login(data) {
  return request({
    url: '/user/login',
    method: 'post',
    data
  })
}

export function getInfo(token) {
  return request({
    url: '/user/info',
    method: 'get',
    params: { token }
  })
}

export function logout() {
  return request({
    url: '/user/logout',
    method: 'post'
  })
}

export function fetchList(query) {
  return request({
    url: '/user/list',
    method: 'get',
    params: query
  })
}

export function fetchDetail(id) {
  return request({
    url: '/user/detail',
    method: 'get',
    params: { id }
  })
}

export function add(data) {
  return request({
    url: '/user/add',
    method: 'post',
    data
  })
}

export function update(data) {
  return request({
    url: '/user/update',
    method: 'post',
    data
  })
}

export function del(data) {
  return request({
    url: `/user/delete`,
    method: 'post',
    data
  })
}

export function addAuth(data) {
  return request({
    url: '/user/addAuth',
    method: 'post',
    data
  })
}

export function delAuth(data) {
  return request({
    url: '/user/delAuth',
    method: 'post',
    data
  })
}
