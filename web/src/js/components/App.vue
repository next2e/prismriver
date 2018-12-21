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

    <div class="linkDiv">
      <p class="linkP"> Queue music - enter a link or go random! </p>
      <div class="linkContainer">
        <form class="linkForm">
          <input id="url" class="linkInput" type="text" autocomplete="off" autofocus placeholder="Insert Youtube url here"></input>
          <button class="linkButton hvr-shuttegit gitr-out-horizontal" id="add-song" type="button">Add</button>
        </form>
        <button id="random" class="hvr-shutter-out-horizontal" type="button">Random</button>
      </div>
    </div>

    <div class="component" style="padding-top: 6px;">
      <p class="linkP"> Choose from past songs &nbsp; <a href="#" id="toggle-search"><span id="search-arrow" class="glyphicon glyphicon-chevron-down"></span></a></p>
      <div id="search-container" style="display: none">
        <div class="linkContainer">
          <form class="linkForm">
            <input id="query" class="linkInput" type="text" autocomplete="off" autofocus placeholder="Search">
            <button id="submit-search" class="linkButton hvr-shutter-out-horizontal" type="button">Search</button>
          </form>
          <button id="shuffle" class="hvr-shutter-out-horizontal" type="button">Shuffle</button>
        </div>
        <ul id="search-results">
        </ul>
      </div>
    </div>

    <div class="component">
      <p class="title"> Currently playing </p>
      <div id="playing">
        <p>Loading...</p>
      </div>
      <div class="skipDiv">
        <button id="skipButton" class="hvr-shutter-out-horizontal"><span class="glyphicon glyphicon-forward"></span>Skip song</button>
      </div>
      <div class="volumeDiv">
        <p class="noMargin">
          Volume:
          <button id="voldown" class="hvr-bounce-to-left"><span class="glyphicon glyphicon-volume-down"></span></button>
          <!-- <span id="vol"> eman </span> -->
          <button id="volup" class="hvr-bounce-to-right"><span class="glyphicon glyphicon-volume-up"></span></button>
        </p>
      </div>
    </div>

    <div class="component">
      <p class="title"> Current queue </p>
      <ul id="queue" class="nomargin">
      </ul>
    </div>
  </div>
</template>