import request from '@/utils/request'

export function fetchList(query) {
  return request({
    url: '/chaincode/list',
    method: 'get',
    params: query
  })
}

export function fetchAllList() {
  return request({
    url: '/chaincode/allList',
    method: 'get'
  })
}

export function fetch(data) {
  return request({
    url: '/chaincode/get',
    method: 'post',
    data
  })
}

export function add(data) {
  return request({
    url: '/chaincode/add',
    method: 'post',
    data
  })
}

export function upgrade(data) {
  return request({
    url: `/chaincode/upgrade`,
    method: 'post',
    data
  })
}

export function deploy(data) {
  return request({
    url: `/chaincode/deploy`,
    method: 'post',
    data
  })
}

export function invoke(data) {
  return request({
    url: `/chaincode/invoke`,
    method: 'post',
    data
  })
}

export function query(data) {
  return request({
    url: `/chaincode/query`,
    method: 'post',
    data
  })
}

export function del(data) {
  return request({
    url: `/chaincode/delete`,
    method: 'post',
    data
  })
}
