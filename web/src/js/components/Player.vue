<template>
  <div class="component">
    <p class="title"> Currently playing </p>
    <div id="playing">
      <p v-if="currenttitle" class="text-overflow">{{ currenttitle }}</p>
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
        <input type="range" :min="0" :max="totalTime / 1000" :value="currentTime / totalTime * totalTime / 1000">
      </p>
    </div>
  </div>
</template>

<script lang="ts">
  import $ from 'jquery'
  import { Component, Prop, Vue } from 'vue-property-decorator'

  @Component
  export default class Player extends Vue {
    currentTime = 0
    state = 0
    totalTime = 0

    @Prop(String) currenttitle!: string

    mounted () {
      setInterval(() => {
        if (this.state === 1) {
          this.currentTime += 1000
        }
      }, 1000)

      const socket = new WebSocket((window.location.protocol === 'https:' ? 'wss://' : 'ws://') +
          window.location.hostname + '/ws/player')

      socket.addEventListener('close', () => {
        alert('WebSocket connection closed. Refresh to continue receiving queue updates!')
      })

      socket.addEventListener('error', () => {
        alert('Error in the WebSocket connection. If there are issues with the queue updating, refresh!')
      })

      socket.addEventListener('message', (event) => {
        const data = JSON.parse(event.data)
        this.currentTime = data.CurrentTime
        this.state = data.State
        this.totalTime = data.TotalTime
      })
    }

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