<template>
  <li>
    <span class="text-overflow">
      <button @click="deleteSong" class="delete"><span class="glyphicon glyphicon-trash"></span></button>
      <button :disabled="disableup" @click="up" class="select-song pull-right"><span class="glyphicon glyphicon-arrow-up"></span></button>
      <button :disabled="disabledown" @click="down" class="select-song pull-right" style="margin-left: 15px;"><span class="glyphicon glyphicon-arrow-down"></span></button>
      {{ item.media.Title }}
    </span>
    <progress-bar v-if="item.downloading" :progress="item.progress"></progress-bar>
  </li>
</template>

<script lang="ts">
  import { Component, Prop, Vue } from 'vue-property-decorator'

  @Component
  export default class QueueItem extends Vue {
    @Prop(Boolean) disabledown!: boolean
    @Prop(Boolean) disableup!: boolean
    @Prop(Number) index!: number
    @Prop(Object) item!: IQueueItem

    deleteSong () {
      fetch('queue/' + this.index, {
        method: 'DELETE'
      })
    }

    down () {
      const params = new URLSearchParams()
      params.set('move', 'down')
      fetch('queue/' + this.index, {
        body: params,
        method: 'PUT'
      })
    }

    up () {
      const params = new URLSearchParams()
      params.set('move', 'up')
      fetch('queue/' + this.index, {
        body: params,
        method: 'PUT'
      })
    }
  }
</script>