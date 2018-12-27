import $ from 'jquery'
import Snackbar from 'node-snackbar'
import Vue, { VNode } from 'vue'

(window as any).jQuery = $

import 'bootstrap'

import App from './components/App.vue'
import QueueItem from './components/QueueItem.vue'
import SearchForm from './components/SearchForm.vue'
import SearchItem from './components/SearchItem.vue'

import '../css/app.css'

function show_msg (text: string): void {
  Snackbar.show({
    duration: 3000,
    showAction: false,
    text
  })
}

function submit_video (): void {
  const url = $('#url')
  $.post('/queue', { url: url.val() })
  url.val('')
  $('#linkButton').blur()
  show_msg('Submitted! Now downloading song...')
}

$(() => {
  Vue.component('queue-item', QueueItem)
  Vue.component('search-form', SearchForm)
  Vue.component('search-item', SearchItem)
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
  $('#random').on('click', () => {
    $('#random').blur()
    show_msg('Adding a random song')
    $.getJSON('/media/random', { limit: 1 }, (data) => {
      $.post('/queue', {
        id: data[0].ID,
        type: data[0].Type
      })
    })
  })
  $('#url').keypress((e) => {
    if (e.keyCode === 13) {
      submit_video()
      e.preventDefault()
    }
  })
})
