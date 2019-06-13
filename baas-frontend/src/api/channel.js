import request from '@/utils/request'

export function fetch(data) {
  return request({
    url: '/channel/get',
    method: 'post',
    data
  })
}

export function fetchAllList(id) {
  return request({
    url: '/channel/allList',
    method: 'get',
    params: { chainId: id }
  })
}

export function add(data) {
  return request({
    url: '/channel/add',
    method: 'post',
    data
  })
}

export function del(data) {
  return request({
    url: `/channel/delete`,
    method: 'post',
    data
  })
}
