<template>
  <div class="component">
    <p class="title"> Currently playing </p>
    <div id="playing">
      <p v-if="item" class="text-overflow">{{ item.media.Title }}</p>
      <p v-else>Nothing currently playing</p>
      <progress-bar v-if="item && item.downloading || item &&  item.error" :progress="item.progress" :error="item.error"></progress-bar>
    </div>
    <div class="skipDiv">
      <button id="skipButton" class="hvr-shutter-out-horizontal" @click="skip"><span class="glyphicon glyphicon-forward"></span>Skip song</button>
    </div>
    <p class="volumeDiv">
      <p class="noMargin">
        Volume:
        <button id="voldown" class="hvr-bounce-to-left" @click="volDown"><span class="glyphicon glyphicon-volume-down"></span></button>
        <!-- <span id="vol"> eman </span> -->
        <span>{{ volume }}</span>
        <button id="volup" class="hvr-bounce-to-right" @click="volUp"><span class="glyphicon glyphicon-volume-up"></span></button>
        <div style="display: flex; width: 100%; align-items: center;">
          <input id="time" type="range" :min="0" :max="totalTime" v-model.number="currentTime"
                 style="width: auto; flex-grow: 2;" @mousedown="seeking = true" @mouseup="seeking = false"
                 @change="seek">
          <span style="margin-left: 15px;">{{ parseTime(currentTime / 1000) +  ' / ' + parseTime(totalTime / 1000) }}</span>
        </div>
      </p>
    </p>
  </div>
</template>

<script lang="ts">
  import $ from 'jquery'
  import { Component, Prop, Vue, Watch } from 'vue-property-decorator'

  @Component
  export default class Player extends Vue {
    currentTime = 0
    seeking = false
    socket: WebSocket | undefined
    state = 0
    totalTime = 0
    volume = 100
    ws = 0

    @Prop(Object) item!: IMedia

    connectWS () {
      this.socket = new WebSocket((window.location.protocol === 'https:' ? 'wss://' : 'ws://') +
          window.location.host + (window.location.pathname === '/' ? '' : window.location.pathname) + '/ws/player')

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

    parseTime (time: number, recur: boolean = false): string {
      time = Math.floor(time)
      let timeString = ''
      timeString = time % 60 < 10 ? '0' + time % 60 + timeString : time % 60 + timeString
      if (time / 60 < 1 && !recur) {
        return '0:' + timeString
      } else if (time / 60 < 1) {
        return timeString
      } else {
        return this.parseTime(time / 60, true) + ':' + timeString
      }
    }

    mounted () {
      this.connectWS()
      setInterval(() => {
        if (this.state === 1 && !this.seeking) {
          this.currentTime += 1000
        }
      }, 1000)

      setInterval(() => {
        if (this.ws === 0) {
          this.connectWS()
        }
      }, 5000)
    }

    @Watch('ws')
    onWSChanged(state: number) {
      this.$emit('update:ws', state)
    }

    seek (event: Event) {
      if (this.state !== 1) {
        return
      }
      const params = new URLSearchParams()
      params.set('seek', (event.target as any).value)
      fetch('player', {
        body: params,
        method: 'PUT'
      })
    }

    skip (event: Event) {
      $((event.target as Object)).blur()
      fetch('queue/0', {
        method: 'DELETE'
      })
    }

    volDown (event: Event) {
      $((event.target as Object)).blur()
      const params = new URLSearchParams()
      params.set('volume', 'down')
      fetch('player', {
        body: params,
        method: 'PUT'
      })
    }

    volUp (event: Event) {
      $((event.target as Object)).blur()
      const params = new URLSearchParams()
      params.set('volume', 'up')
      fetch('player', {
        body: params,
        method: 'PUT'
      })
    }
  }
</script>
