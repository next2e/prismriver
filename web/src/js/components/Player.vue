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
        <span>{{ volume }}</span>
        <button id="volup" class="hvr-bounce-to-right" @click="volUp"><span class="glyphicon glyphicon-volume-up"></span></button>
        <input type="range" :min="0" :max="totalTime / 1000" :value="currentTime / totalTime * totalTime / 1000">
      </p>
    </div>
  </div>
</template>

<script lang="ts">
  import $ from 'jquery'
  import { Component, Prop, Vue, Watch } from 'vue-property-decorator'

  @Component
  export default class Player extends Vue {
    currentTime = 0
    socket: WebSocket | undefined
    state = 0
    totalTime = 0
    volume = 100
    ws = 0

    @Prop(String) currenttitle!: string

    mounted () {
      setInterval(() => {
        if (this.state === 1) {
          this.currentTime += 1000
        }
      }, 1000)

      setInterval(() => {
        if (this.ws === 0) {
          this.socket = new WebSocket((window.location.protocol === 'https:' ? 'wss://' : 'ws://') +
              window.location.hostname + '/ws/player')

          this.socket.addEventListener('close', () => {
            this.ws = 0
          })
          this.socket.addEventListener('error', () => {
            this.ws = 2
          })
          this.socket.addEventListener('message', (event) => {
            this.ws = 1
            const data = JSON.parse(event.data)
            this.currentTime = data.CurrentTime
            this.state = data.State
            this.totalTime = data.TotalTime
            this.volume = data.Volume
          })
        }
      }, 5000)
    }

    @Watch('ws')
    onWSChanged(state: number) {
      this.$emit('update:ws', state)
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