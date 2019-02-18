<template>
  <div class="linkDiv">
    <p class="linkP"> Queue music - enter a link or go random! </p>
    <div class="linkContainer">
      <form class="linkForm">
        <input id="url" class="linkInput" v-model="url" type="text" autocomplete="off" autofocus placeholder="Insert Youtube url here" @keydown.enter.prevent="submit">
        <button class="linkButton hvr-shuttegit gitr-out-horizontal" id="add-song" type="button" @click="submit">Add</button>
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
    url = ''

    random (event: Event) {
      $((event.target as Object)).blur()
      this.showMessage('Adding a random song')
      $.getJSON(window.location.toString() + '/media/random', { limit: 1 }, (data: [{ ID: number, Type: string }]) => {
        $.post(window.location.toString() + '/queue', {
          id: data[0].ID,
          type: data[0].Type
        })
      })
    }

    submit () {
      $.post(window.location.toString() + '/queue', { url: this.url })
      this.showMessage('Submitted! Now downloading song...')
      this.url = ''
    }
  }
</script>