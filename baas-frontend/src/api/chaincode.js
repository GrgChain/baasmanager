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

export function queryLedger(id) {
  return request({
    url: `/chaincode/queryLedger?channelId=` + id,
    method: 'get'
  })
}

export function queryLatestBlocks(id) {
  return request({
    url: `/chaincode/queryLatestBlocks?channelId=` + id,
    method: 'get'
  })
}

export function queryBlock(cid, search) {
  return request({
    url: `/chaincode/queryBlock?channelId=` + cid + `&search=` + search,
    method: 'get'
  })
}

export function del(data) {
  return request({
    url: `/chaincode/delete`,
    method: 'post',
    data
  })
}
