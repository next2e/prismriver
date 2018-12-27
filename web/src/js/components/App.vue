<template>
  <div id="main">
    <!-- <span>Golang give me the power to fix this CSS</span> -->
    <header>
      <div class="container">
        <div class="first">
          <span>2GroovE</span>
          <i class="glyphicon glyphicon-music"></i>
        </div>
        <button class="hvr-pulse" id="quietButton"><i class="glyphicon glyphicon-exclamation-sign"></i><span class="quietText">Be quiet!</span></button>
      </div>
    </header>

    <url-form></url-form>
    <search-form></search-form>

    <div class="component">
      <p class="title"> Currently playing </p>
      <div id="playing">
        <p v-if="queue.length">{{ queue[0] }}</p>
        <p v-else>Nothing currently playing</p>
      </div>
      <div class="skipDiv">
        <button id="skipButton" class="hvr-shutter-out-horizontal" @click="skip"><span class="glyphicon glyphicon-forward"></span>Skip song</button>
      </div>
      <div class="volumeDiv">
        <p class="noMargin">
          Volume:
          <button id="voldown" class="hvr-bounce-to-left" @click="volDown"><span class="glyphicon glyphicon-volume-down"></span></button>
          <!-- <span id="vol"> eman </span> -->
          <button id="volup" class="hvr-bounce-to-right" @click="volUp"><span class="glyphicon glyphicon-volume-up"></span></button>
        </p>
      </div>
    </div>

    <div class="component">
      <p class="title"> Current queue </p>
      <ul id="queue" class="nomargin">
        <queue-item v-for="(item, index) in queue" v-if="index > 0" :index="index" :title="item"></queue-item>
      </ul>
    </div>
  </div>
</template>

<script lang="ts">
  import $ from 'jquery'
  import Vue from 'vue'
  import Component from 'vue-class-component'

  const BaseComponent = Vue.extend({
    data () {
      return {
        queue: [],
        results: []
      }
    },
    mounted () {
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
        this.queue = queue
      })
    }
  })

  @Component
  export default class App extends BaseComponent {
    skip (event: Event) {
      $((event.target as Object)).blur()
      $.ajax({
        type: 'DELETE',
        url: '/queue/0'
      })
    }

    volDown (event: Event) {
      $((event.target as Object)).blur()
      $.ajax({
        data: { volume: 'down' },
        type: 'PUT',
        url: '/player'
      })
    }

    volUp (event: Event) {
      $((event.target as Object)).blur()
      $.ajax({
        data: { volume: 'up' },
        type: 'PUT',
        url: '/player'
      })
    }
  }
</script>