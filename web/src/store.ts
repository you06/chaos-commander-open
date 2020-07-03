import Vue from 'vue'
import Vuex from 'vuex'
import axios from 'axios'

Vue.use(Vuex)

const instance = axios.create({
  baseURL: 'http://localhost:50000/api/v1/',
  timeout: 30000,
  headers: {},
})

export default new Vuex.Store({
  state: {
    jobs: []
  },
  mutations: {

  },
  actions: {
    getJobs(ctx) {
      instance({
        method: 'get',
        url: '/jobs',
      }).then((res) => {
        console.log(res)
      }).catch((reason) => {
        console.log(reason)
      })
    },
  },
})
