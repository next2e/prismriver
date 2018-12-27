import $ from 'jquery'
import Snackbar from 'node-snackbar'
import Vue, { VNode } from 'vue'

(window as any).jQuery = $

import 'bootstrap'

import App from './components/App.vue'
import QueueItem from './components/QueueItem.vue'
import SearchForm from './components/SearchForm.vue'
import SearchItem from './components/SearchItem.vue'
import URLForm from './components/URLForm.vue'

import '../css/app.css'

$(() => {
  Vue.component('queue-item', QueueItem)
  Vue.component('search-form', SearchForm)
  Vue.component('search-item', SearchItem)
  Vue.component('url-form', URLForm)
  Vue.component('app', App)

  Vue.mixin({
    methods: {
      showMessage (message: string): void {
        Snackbar.show({
          duration: 3000,
          showAction: false,
          text: message
        })
      }
    }
  }); // not sure why this is needed here but TS is whining

  (window as any).app = new Vue({
    el: '#main',
    render (createElement: (el: string) => VNode): VNode {
      return createElement('app')
    }
  })
  $('#quietButton').on('click', () => {
    $('#quietButton').blur()
    $.get('/quiet')
  })
})
