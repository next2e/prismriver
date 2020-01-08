<template>
  <div class="linkDiv">
    <p class="linkP"> Queue music - enter a link or go random! </p>
    <div class="linkContainer">
      <form class="linkForm" style="flex-direction: column; align-items: flex-start;">
        <div style="display: flex; width: 100%;">
          <input id="url" class="linkInput" v-model="url" type="text" autocomplete="off" autofocus placeholder="Insert Youtube url here" @keydown.enter.prevent="submit">
          <button class="linkButton hvr-shuttegit gitr-out-horizontal" id="add-song" type="button" @click="submit">Add</button>
        </div>
        <p><input type="checkbox" id="video" v-model="video"><label for="video">With Video?</label></p>
      </form>
      <button id="random" class="hvr-shutter-out-horizontal" type="button" @click="random">Random</button>
    </div>
  </div>
</template>

<script lang="ts">
  import $ from 'jquery'
  import { Component, Mixins } from 'vue-property-decorator'

  import ShowMessage from '../mixins/ShowMessage'

  @Component
  export default class URLForm extends Mixins(ShowMessage) {
    video = false
    url = ''

    random (event: Event) {
      $((event.target as Object)).blur()
      this.showMessage('Adding a random song')
      const params = new URLSearchParams()
      params.set('limit', '1')
      fetch('media/random?' + params.toString()).then((response) => {
        return response.json()
      }).then((json) => {
        const data = json[0]
        const params = new URLSearchParams()
        params.set('id', data.ID)
        params.set('type', data.Type)
        return fetch('queue', {
          body: params,
          method: 'POST'
        })
      })
    }

    submit () {
      const params = new URLSearchParams()
      params.set('url', this.url)
      params.set('video', this.video.toString())
      fetch('queue', {
        body: params,
        method: 'POST'
      })
      this.showMessage('Submitted! Now downloading song...')
      this.url = ''
      this.video = false
    }
  }
</script>