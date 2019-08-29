import request from '@/utils/request'

export function counts(query) {
  return request({
    url: '/dashboard/counts',
    method: 'get',
    params: query
  })
}

export function consensusTotal(query) {
  return request({
    url: '/dashboard/consensusTotal',
    method: 'get',
    params: query
  })
}

export function getSevenDays(query) {
  return request({
    url: '/dashboard/sevenDays',
    method: 'get',
    params: query
  })
}
