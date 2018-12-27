import $ from 'jquery'
import Snackbar from 'node-snackbar'
import Vue, { VNode } from 'vue'

(window as any).jQuery = $

import 'bootstrap'

import App from './components/App.vue'
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

(window as any).delete_song = (songNumber: number): void => {
  $.ajax({
    type: 'DELETE',
    url: '/queue/' + songNumber
  })
}

$(() => {
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
  $('#volup').on('click', () => {
    $('#volup').blur()
    $.ajax({
      data: { volume: 'up' },
      type: 'PUT',
      url: '/player'
    })
  })
  $('#voldown').on('click', () => {
    $('#voldown').blur()
    $.ajax({
      data: { volume: 'down' },
      type: 'PUT',
      url: '/player'
    })
  })
  $('#skipButton').on('click', () => {
    $('#skipButton').blur()
    $.ajax({
      type: 'DELETE',
      url: '/queue/0'
    })
  })
  $('#url').keypress((e) => {
    if (e.keyCode === 13) {
      submit_video()
      e.preventDefault()
    }
  })

  const socket = new WebSocket((window.location.protocol === 'https:' ? 'wss://' : 'ws://') +
      window.location.hostname + '/ws/queue')

  socket.addEventListener('close', () => {
    alert('WebSocket connection closed. Refresh to continue receiving queue updates!')
  })

  socket.addEventListener('error', () => {
    alert('Error in the WebSocket connection. If there are issues with the queue updating, refresh!')
  })

  socket.addEventListener('message', (event) => {
    const queue = JSON.parse(event.data)
    let things = ''
    $('#playing').html('<p>' + queue[0] + '</p>')
    const upcoming = queue.slice(1)
    for (const index in upcoming) {
      if (!upcoming.hasOwnProperty(index)) {
        continue
      }
      things += '<li><span>'
      things += '<button onclick="delete_song(' + (parseInt(index, 10) + 1) + ')" class="delete"> Delete </button>'
      things += upcoming[index] + '</span></li>'
    }
    $('#queue').html(things)
  })
})
