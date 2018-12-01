import $ from 'jquery'
import Snackbar from 'node-snackbar'

import '../css/app.css'

function show_msg (text: string): void {
  Snackbar.show({
    duration: 3000,
    showAction: false,
    text
  })
}

// Sanitize punctuation
function encode_str (s: string): string {
  return s.replace(/'/g, '%27').replace(/"/g, '%22')
}

function decode_str (s: string): string {
  return s.replace(/%27/g, "'").replace(/%22/g, '"')
}

function submit_video (): void {
  const url = $('#url')
  $.post('/queue', { url: url.val() })
  url.val('')
  $('#linkButton').blur()
  show_msg('Submitted! Now downloading song...')
}

function submit_search (): void {
  const query = $('#query').val()
  if ((query as string).length) {
    $.getJSON('/media/search', { query }, display_results)
  } else {
    $.getJSON('/media/random', { limit: 20 }, display_results)
  }
}

(window as any).delete_song = (songNumber: number): void => {
  $.ajax({
    type: 'DELETE',
    url: '/queue/' + songNumber
  })
}

(window as any).select_song = (song: string, type: string, title: string): void => {
  title = decode_str(title)
  show_msg('Adding ' + title)
  $.post('/queue', { id: song, type })
}

function display_results (songs: string[]): void {
  let listing = ''
  songs.forEach((song) => {
    listing += '<li></span><button onclick="select_song(\'' + (song as any).ID + '\', \'' + (song as any).Type +
        '\', \'' + encode_str((song as any).Title) + '\')" '
    listing += 'class="select-song">Add</button> '
    listing += (song as any).Title
    listing += '</span></li>'
  })
  $('#search-results').html(listing)
}

$(() => {
  $.getJSON('/media/random', { limit: 20 }, display_results)
  $('#add-song').on('click', submit_video)
  $('#submit-search').on('click', submit_search)
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
  $('#shuffle').on('click', () => {
    $('#shuffle').blur()
    $.getJSON('/media/random', { limit: 20 }, display_results)
  })
  $('#toggle-search').on('click', () => {
    if ($('#search-container:visible').length) {
      $('#search-container').hide('slow')
    } else {
      $('#search-container').show('slow')
    }
    const searchArrow = $('#search-arrow')
    searchArrow.toggleClass('glyphicon-chevron-down')
    searchArrow.toggleClass('glyphicon-chevron-up')
    $('#toggle-search').blur()
  })
  $('#url').keypress((e) => {
    if (e.keyCode === 13) {
      submit_video()
      e.preventDefault()
    }
  })
  $('#query').keypress((e) => {
    if (e.keyCode === 13) {
      submit_search()
      e.preventDefault()
    }
  })

  const socket = new WebSocket('ws://' + window.location.hostname + '/ws/queue')

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
