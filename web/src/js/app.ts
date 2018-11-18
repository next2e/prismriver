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
  const url = $('url')
  $.post('/submit', { video: url.val() })
  url.val('')
  $('#linkButton').blur()
  show_msg('Submitted! Now downloading song...')
}

function submit_search (): void {
  const query = $('#query').val()
  if ((query as string).length) {
    $.getJSON('/search', { query }, display_results)
  } else {
    $.getJSON('/shuffle', display_results)
  }
}

function delete_song (songNumber: number): void {
  $.post('/delete', { song: songNumber })
}

function select_song (song: string): void {
  song = decode_str(song)
  const title = song.substring(0, song.lastIndexOf('.'))
  show_msg('Adding ' + title)
  $.post('/playsong', { song })
}

function display_results (songs: string[]): void {
  let listing = ''
  songs.forEach((song) => {
    listing += '<li></span><button onclick="select_song(\'' + encode_str(song) + '\')" '
    listing += 'class="select-song">Add</button> '
    listing += song.substring(0, song.lastIndexOf('.'))
    listing += '</span></li>'
  })
  $('#search-results').html(listing)
}

$(() => {
  $.getJSON('/shuffle', display_results)
  $('#add-song').on('click', submit_video)
  $('#submit-search').on('click', submit_search)
  $('#quietButton').on('click', () => {
    $('#quietButton').blur()
    $.get('/quiet')
  })
  $('#random').on('click', () => {
    $('#random').blur()
    show_msg('Adding a random song')
    $.get('/random')
  })
  $('#volup').on('click', () => {
    $('#volup').blur()
    $.get('/volup')
  })
  $('#voldown').on('click', () => {
    $('#voldown').blur()
    $.get('/voldown')
  })
  $('#skipButton').on('click', () => {
    $('#skipButton').blur()
    $.get('/skip')
  })
  $('#shuffle').on('click', () => {
    $('#shuffle').blur()
    $.getJSON('/shuffle', display_results)
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
      things += '<button onclick="delete_song(' + index + ')" class="delete"> Delete </button>'
      things += upcoming[index] + '</span></li>'
    }
    $('#queue').html(things)
  })
})
