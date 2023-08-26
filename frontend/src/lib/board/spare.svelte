<script lang="ts">
  import JuicerPiece from './piece.svelte';
  import { Piece } from '$lib/board/board';

  let blackSparePieces: Piece[] = [
    new Piece('r', 'b'),
    new Piece('n', 'b'),
    new Piece('b', 'b'),
    new Piece('q', 'b'),
    new Piece('k', 'b'),
    new Piece('p', 'b'),
  ];

  let whiteSparePieces: Piece[] = [
    new Piece('r', 'w'),
    new Piece('n', 'w'),
    new Piece('b', 'w'),
    new Piece('q', 'w'),
    new Piece('k', 'w'),
    new Piece('p', 'w'),
  ];

  let sparePieces = blackSparePieces.concat(whiteSparePieces);
</script>

<div class="spare-pieces">
  {#each blackSparePieces as piece, idx}
    <JuicerPiece {piece} squareIdx={idx} />
  {/each}
</div>

<style>
  .spare-pieces {
    display: grid;
    grid-template-columns: 6;
    gap: 1rem;
  }
</style>
