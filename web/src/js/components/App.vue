<template>
  <div id="main">
    <!-- <span>Golang give me the power to fix this CSS</span> -->
    <header style="position: relative; z-index: 420;">
      <div class="container">
        <div class="first">
          <span>2GroovE</span>
          <i class="glyphicon glyphicon-music"></i>
        </div>
        <button class="hvr-pulse" id="quietButton" @click="beQuiet"><i class="glyphicon glyphicon-exclamation-sign"></i><span class="quietText">Be quiet!</span></button>
      </div>
    </header>

    <transition name="ws">
      <div v-if="state() !== 1" class="alert" :class="state() === 2 ? 'alert-warning' : 'alert-danger'" style="height: 50px; position: relative;">
        <p v-if="state() === 0">The player has disconnected from the server. Attempting to reconnect...</p>
        <p v-if="state() === 2">The page has encountered errors. The current queue may be outdated.</p>
      </div>
    </transition>

    <url-form></url-form>
    <search-form></search-form>
    <player :item="queue[0]" @update:ws="playerWS = $event"></player>

    <div class="component">
      <p class="title">
        Current queue <button class="select-song" @click="shuffle"><span class="glyphicon glyphicon-random"></span></button>
      </p>
      <transition-group id="queue" class="nomargin" name="queue" tag="ul">
        <queue-item class="queue-item" v-for="(item, index) in queue" :key="item.media.Title" v-if="index > 0" :index="index" :item="item" :disableup="index === 1" :disabledown="index === queue.length - 1"></queue-item>
      </transition-group>
    </div>
  </div>
</template>

<script lang="ts">
  import { Component, Vue } from 'vue-property-decorator'

  @Component
  export default class App extends Vue {
    playerWS = 0
    queue = []
    queueWS = 0
    results = []
    socket: WebSocket | undefined

    beQuiet () {
      const params = new URLSearchParams()
      params.set('quiet', 'true')
      fetch('player', {
        body: params,
        method: 'PUT'
      })
    }

    connectWS() {
      this.socket = new WebSocket((window.location.protocol === 'https:' ? 'wss://' : 'ws://') +
          window.location.host + (window.location.pathname === '/' ? '' : window.location.pathname) + '/ws/queue')

      this.socket.addEventListener('close', () => {
        this.queueWS = 0
      })
      this.socket.addEventListener('error', () => {
        this.queueWS = 2
      })
      this.socket.addEventListener('message', (event) => {
        this.queueWS = 1
        const queue = JSON.parse(event.data)
        this.queue = queue
      })
    }

    mounted () {
      this.connectWS()
      setInterval(() => {
        if (this.playerWS === 0) {
          this.connectWS()
        }
      }, 5000)
    }

    shuffle () {
      const params = new URLSearchParams()
      params.set('shuffle', 'true')
      fetch('player', {
        body: params,
        method: 'PUT'
      })
    }

    state () {
      return Math.max(this.playerWS, this.queueWS)
    }
  }
</script>

<style scoped>
  .queue-enter, .queue-leave-to {
    opacity: 0;
    transform: translateX(30px);
  }
  .queue-item {
    transition: all 1s;
  }
  .queue-leave-active {
    position: absolute;
  }
</style>
