<script>
  import { onMount } from "svelte";

  export let src;
  export let alt = "";
  export let width;

  let imageRef;
  let isVisible = false;

  onMount(() => {
    const observer = new IntersectionObserver(([entry]) => {
      if (entry.isIntersecting) {
        isVisible = true;
        observer.disconnect(); // Stop observing once loaded
      }
    });
    observer.observe(imageRef);
  });
</script>

<img bind:this={imageRef} {alt} {width} src={isVisible ? src : ""} loading="lazy" />
