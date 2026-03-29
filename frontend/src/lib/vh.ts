/**
 * Installs a `--vh` CSS custom property set to 1/100th of the visual viewport height.
 * More reliable than `100vh` on mobile browsers (Android Chrome, iOS Safari) because
 * it updates when the browser UI (address bar, bottom toolbar) resizes the viewport.
 *
 * Usage in CSS: height: calc(var(--vh, 1vh) * 100);
 *
 * @returns Cleanup function that removes the event listeners.
 */
export function installVhVar(): () => void {
  const set = () => {
    const h = window.visualViewport?.height ?? window.innerHeight
    document.documentElement.style.setProperty('--vh', `${h * 0.01}px`)
  }
  set()
  window.visualViewport?.addEventListener('resize', set)
  window.addEventListener('resize', set)
  return () => {
    window.visualViewport?.removeEventListener('resize', set)
    window.removeEventListener('resize', set)
  }
}
