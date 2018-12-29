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
    <player :currenttitle="queue.length ? queue[0].Media.Title : null"></player>

    <div class="component">
      <p class="title"> Current queue </p>
      <ul id="queue" class="nomargin">
        <queue-item v-for="(item, index) in queue" v-if="index > 0" :index="index" :item="item" :disableup="index === 1" :disabledown="index === queue.length - 1"></queue-item>
      </ul>
    </div>
  </div>
</template>

<script lang="ts">
  import { Component, Vue } from 'vue-property-decorator'
  @Component
  export default class App extends Vue {
    queue = []
    results = []

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
  }
</script>