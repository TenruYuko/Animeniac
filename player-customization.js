// This script will be injected into the page to customize the player UI
function customizePlayer() {
  console.log("Player customization script loaded");

  // Function to create our custom player controls
  function createCustomPlayerControls() {
    // Check if the player exists in the DOM
    const playerContainer = document.querySelector('[data-sea-media-player]');
    if (!playerContainer) return;
    
    console.log("Player found, adding custom controls");
    
    // Remove any existing custom controls
    const existingControls = document.getElementById('modern-player-controls');
    if (existingControls) existingControls.remove();
    
    // Create our modern player controls container
    const controlsContainer = document.createElement('div');
    controlsContainer.id = 'modern-player-controls';
    controlsContainer.className = 'player-custom-controls';
    
    // Create rewind controls group
    const rewindGroup = document.createElement('div');
    rewindGroup.className = 'player-rewind-group';
    
    // Create play/pause container (center)
    const centerControls = document.createElement('div');
    centerControls.className = 'player-center-group';
    
    // Create forward controls group
    const forwardGroup = document.createElement('div');
    forwardGroup.className = 'player-forward-group';
    
    // Add rewind buttons (10s, 30s, 2x, 5x, 10x, 20x)
    const rewind10Button = createControlButton('⟲10', () => seekRelative(-10));
    const rewind30Button = createControlButton('⟲30', () => seekRelative(-30));
    const rewind2xButton = createControlButton('⟲2x', () => toggleRewind(2));
    const rewind5xButton = createControlButton('⟲5x', () => toggleRewind(5));
    const rewind10xButton = createControlButton('⟲10x', () => toggleRewind(10));
    const rewind20xButton = createControlButton('⟲20x', () => toggleRewind(20));
    
    // Add forward buttons (10s, 30s, 2x, 5x, 10x, 20x)
    const forward10Button = createControlButton('10⟳', () => seekRelative(10));
    const forward30Button = createControlButton('30⟳', () => seekRelative(30));
    const forward2xButton = createControlButton('2x⟳', () => toggleFastForward(2));
    const forward5xButton = createControlButton('5x⟳', () => toggleFastForward(5));
    const forward10xButton = createControlButton('10x⟳', () => toggleFastForward(10));
    const forward20xButton = createControlButton('20x⟳', () => toggleFastForward(20));
    
    // Add buttons to groups
    rewindGroup.appendChild(rewind20xButton);
    rewindGroup.appendChild(rewind10xButton);
    rewindGroup.appendChild(rewind5xButton);
    rewindGroup.appendChild(rewind2xButton);
    rewindGroup.appendChild(rewind30Button);
    rewindGroup.appendChild(rewind10Button);
    
    forwardGroup.appendChild(forward10Button);
    forwardGroup.appendChild(forward30Button);
    forwardGroup.appendChild(forward2xButton);
    forwardGroup.appendChild(forward5xButton);
    forwardGroup.appendChild(forward10xButton);
    forwardGroup.appendChild(forward20xButton);
    
    // Add groups to container
    controlsContainer.appendChild(rewindGroup);
    controlsContainer.appendChild(centerControls);
    controlsContainer.appendChild(forwardGroup);
    
    // Add to the player
    playerContainer.appendChild(controlsContainer);
    
    // Add our custom CSS
    const style = document.createElement('style');
    style.textContent = `
      /* Custom player controls */
      .vds-play-button {
        position: absolute !important;
        bottom: 20px !important;
        left: 50% !important;
        transform: translateX(-50%) !important;
        z-index: 100 !important;
        height: 60px !important;
        width: 60px !important;
        background-color: rgba(255, 255, 255, 0.2) !important;
        border-radius: 50% !important;
        display: flex !important;
        align-items: center !important;
        justify-content: center !important;
      }

      .vds-play-button:hover {
        background-color: rgba(255, 255, 255, 0.3) !important;
      }

      .vds-play-button svg {
        height: 30px !important;
        width: 30px !important;
      }

      /* Add rewind/forward buttons */
      .player-custom-controls {
        position: absolute !important;
        bottom: 20px !important;
        width: 100% !important;
        display: flex !important;
        justify-content: center !important;
        align-items: center !important;
        z-index: 99 !important;
      }

      .player-rewind-group {
        position: absolute !important;
        left: 20% !important;
        display: flex !important;
        gap: 8px !important;
      }

      .player-forward-group {
        position: absolute !important;
        right: 20% !important;
        display: flex !important;
        gap: 8px !important;
      }

      .player-control-button {
        height: 40px !important;
        width: 40px !important;
        background-color: rgba(255, 255, 255, 0.15) !important;
        border-radius: 50% !important;
        display: flex !important;
        align-items: center !important;
        justify-content: center !important;
        cursor: pointer !important;
        color: white !important;
        font-size: 12px !important;
        font-weight: bold !important;
      }

      .player-control-button:hover {
        background-color: rgba(255, 255, 255, 0.25) !important;
      }

      .player-control-button.active {
        background-color: rgba(255, 255, 255, 0.35) !important;
      }
    `;
    document.head.appendChild(style);
  }
  
  // Helper function to create a control button
  function createControlButton(text, onClick) {
    const button = document.createElement('button');
    button.className = 'player-control-button';
    button.textContent = text;
    button.onclick = onClick;
    return button;
  }
  
  // Function to seek relative to current position
  function seekRelative(seconds) {
    const player = document.querySelector('[data-sea-media-player]');
    if (!player || !player.currentTime) return;
    
    const currentTime = player.currentTime;
    const duration = player.duration || 0;
    const newTime = Math.max(0, Math.min(currentTime + seconds, duration));
    player.currentTime = newTime;
  }
  
  // Variables to track fast-forward and rewind state
  let isFastForwarding = false;
  let isRewinding = false;
  let ffSpeed = 2;
  let rwSpeed = 2;
  let ffInterval = null;
  let rwInterval = null;
  
  // Function to toggle fast-forward
  function toggleFastForward(speed) {
    if (isFastForwarding && speed === ffSpeed) {
      // Turn off fast-forward
      clearInterval(ffInterval);
      ffInterval = null;
      isFastForwarding = false;
    } else {
      // Turn off rewind if it's on
      if (isRewinding) {
        clearInterval(rwInterval);
        rwInterval = null;
        isRewinding = false;
      }
      
      // Update speed
      if (speed) ffSpeed = speed;
      
      // Clear existing interval
      if (ffInterval) {
        clearInterval(ffInterval);
      }
      
      // Start new interval
      ffInterval = setInterval(() => {
        seekRelative(ffSpeed);
      }, 1000);
      
      isFastForwarding = true;
    }
    
    // Update control visuals
    updateControlsVisuals();
  }
  
  // Function to toggle rewind
  function toggleRewind(speed) {
    if (isRewinding && speed === rwSpeed) {
      // Turn off rewind
      clearInterval(rwInterval);
      rwInterval = null;
      isRewinding = false;
    } else {
      // Turn off fast-forward if it's on
      if (isFastForwarding) {
        clearInterval(ffInterval);
        ffInterval = null;
        isFastForwarding = false;
      }
      
      // Update speed
      if (speed) rwSpeed = speed;
      
      // Clear existing interval
      if (rwInterval) {
        clearInterval(rwInterval);
      }
      
      // Start new interval
      rwInterval = setInterval(() => {
        seekRelative(-rwSpeed);
      }, 1000);
      
      isRewinding = true;
    }
    
    // Update control visuals
    updateControlsVisuals();
  }
  
  // Function to update the visuals of active/inactive controls
  function updateControlsVisuals() {
    // Update in next tick
    setTimeout(() => {
      // Find and update all control buttons
      const rewindButtons = document.querySelectorAll('.player-control-button');
      rewindButtons.forEach(button => {
        const text = button.textContent;
        
        // Reset all buttons first
        button.classList.remove('active');
        
        // Set active state based on current state
        if (isRewinding && text.includes(`⟲${rwSpeed}x`)) {
          button.classList.add('active');
        }
        
        if (isFastForwarding && text.includes(`${ffSpeed}x⟳`)) {
          button.classList.add('active');
        }
      });
    }, 0);
  }
  
  // Run on initial load and whenever the URL changes
  createCustomPlayerControls();

  // Observe DOM changes to detect when the player is added
  const observer = new MutationObserver((mutations) => {
    for (const mutation of mutations) {
      if (mutation.addedNodes.length) {
        createCustomPlayerControls();
      }
    }
  });
  
  // Start observing the document body
  observer.observe(document.body, { childList: true, subtree: true });
  
  // Add keyboard shortcuts
  document.addEventListener('keydown', (event) => {
    const key = event.key.toLowerCase();
    
    // Basic shortcuts
    if (key === 'arrowright') seekRelative(10);
    if (key === 'arrowleft') seekRelative(-10);
    
    // Shift + arrows for 30 seconds
    if (event.shiftKey && key === 'arrowright') seekRelative(30);
    if (event.shiftKey && key === 'arrowleft') seekRelative(-30);
    
    // Alt + number for speeds
    if (event.altKey) {
      const num = parseInt(key, 10);
      if (num === 1) toggleRewind(2);
      if (num === 2) toggleRewind(5);
      if (num === 3) toggleRewind(10);
      if (num === 4) toggleRewind(20);
      if (num === 5) toggleFastForward(2);
      if (num === 6) toggleFastForward(5);
      if (num === 7) toggleFastForward(10);
      if (num === 8) toggleFastForward(20);
    }
  });
  
  // Clean up on page unload
  window.addEventListener('beforeunload', () => {
    if (ffInterval) clearInterval(ffInterval);
    if (rwInterval) clearInterval(rwInterval);
    observer.disconnect();
  });
}

// Insert script tag to run when page loads
document.addEventListener('DOMContentLoaded', customizePlayer);
