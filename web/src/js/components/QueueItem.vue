<template>
  <li>
    <span class="text-overflow">
      <button @click="deleteSong" class="delete"><span class="glyphicon glyphicon-trash"></span></button>
      <button :disabled="disableup" @click="up" class="select-song pull-right"><span class="glyphicon glyphicon-arrow-up"></span></button>
      <button :disabled="disabledown" @click="down" class="select-song pull-right" style="margin-left: 15px;"><span class="glyphicon glyphicon-arrow-down"></span></button>
      {{ item.Media.Title }}
    </span>
    <progress-bar v-if="item.Downloading" :progress="item.DownloadProgress"></progress-bar>
  </li>
</template>

<script lang="ts">
  import $ from 'jquery'
  import { Component, Prop, Vue } from 'vue-property-decorator'

  @Component
  export default class QueueItem extends Vue {
    @Prop(Boolean) disabledown!: boolean
    @Prop(Boolean) disableup!: boolean
    @Prop(Number) index!: number
    @Prop(Object) item!: IQueueItem

    deleteSong () {
      $.ajax({
        type: 'DELETE',
        url: window.location.toString() + '/queue/' + this.index
      })
    }

    down () {
      $.ajax({
        data: { move: 'down' },
        type: 'PUT',
        url: window.location.toString() + '/queue/' + this.index
      })
    }

    up () {
      $.ajax({
        data: { move: 'up' },
        type: 'PUT',
        url:window.location.toString() + '/queue/' + this.index
      })
    }
  }
</script>