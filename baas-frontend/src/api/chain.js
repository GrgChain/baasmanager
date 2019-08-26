import request from '@/utils/request'

export function fetchList(query) {
  return request({
    url: '/chain/list',
    method: 'get',
    params: query
  })
}

export function fetchAllList() {
  return request({
    url: '/chain/allList',
    method: 'get'
  })
}

export function fetch(data) {
  return request({
    url: '/chain/get',
    method: 'post',
    data
  })
}

export function add(data) {
  return request({
    url: '/chain/add',
    method: 'post',
    data
  })
}

export function update(data) {
  return request({
    url: `/chain/update`,
    method: 'post',
    data
  })
}

export function del(data) {
  return request({
    url: `/chain/delete`,
    method: 'post',
    data
  })
}

export function build(data) {
  return request({
    url: `/chain/build`,
    method: 'post',
    data
  })
}

export function run(data) {
  return request({
    url: `/chain/run`,
    method: 'post',
    data
  })
}

export function stop(data) {
  return request({
    url: `/chain/stop`,
    method: 'post',
    data
  })
}

export function release(data) {
  return request({
    url: `/chain/release`,
    method: 'post',
    data
  })
}

export function podsQuery(id) {
  return request({
    url: `/chain/podsQuery?chainId=` + id,
    method: 'get'
  })
}

export function changeSize(data) {
  return request({
    url: `/chain/changeSize`,
    method: 'post',
    data
  })
}

export function download(id) {
  const aTag = document.createElement('a')
  aTag.href = process.env.VUE_APP_BASE_API + '/chain/download?chainId=' + id
  aTag.click()
  return true
}
