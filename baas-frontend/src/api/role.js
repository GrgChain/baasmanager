import request from '@/utils/request'

export function getRoutes() {
  return request({
    url: '/routes',
    method: 'get'
  })
}

export function fetchList(query) {
  return request({
    url: '/role/list',
    method: 'get',
    params: query
  })
}

export function fetchAllList() {
  return request({
    url: '/role/allList',
    method: 'get'
  })
}

export function add(data) {
  return request({
    url: '/role/add',
    method: 'post',
    data
  })
}

export function update(data) {
  return request({
    url: `/role/update`,
    method: 'post',
    data
  })
}

export function del(data) {
  return request({
    url: `/role/delete`,
    method: 'post',
    data
  })
}
