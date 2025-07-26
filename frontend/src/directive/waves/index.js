import './waves.css'

const context = '@@wavesContext'

function isRippleEnabled(value) {
  return value !== false
}

function rippleShow(e) {
  const container = e.currentTarget
  const instance = container[context]

  // Create ripple
  const ripple = document.createElement('div')
  ripple.className = 'waves-ripple'

  // Get click coordinates
  const rect = container.getBoundingClientRect()
  const left = rect.left
  const top = rect.top
  const width = container.offsetWidth
  const height = container.offsetHeight
  const dx = e.clientX - left
  const dy = e.clientY - top
  const maxX = Math.max(dx, width - dx)
  const maxY = Math.max(dy, height - dy)
  const style = window.getComputedStyle(container)
  const radius = Math.sqrt(maxX * maxX + maxY * maxY)

  // Apply styles
  ripple.style.top = dy + 'px'
  ripple.style.left = dx + 'px'
  ripple.style.width = radius * 2 + 'px'
  ripple.style.height = radius * 2 + 'px'
  ripple.style.marginLeft = -radius + 'px'
  ripple.style.marginTop = -radius + 'px'

  // Remove previous ripples
  container.querySelectorAll('.waves-ripple').forEach(el => el.remove())

  // Add ripple to container
  container.appendChild(ripple)

  // Animate
  ripple.classList.add('waves-ripple-animate')
}

function rippleHide(e) {
  const container = e.currentTarget
  const ripples = container.querySelectorAll('.waves-ripple')

  if (ripples.length === 0) return

  const ripple = ripples[ripples.length - 1]
  const delay = 0.3

  const duration = parseFloat(
    window.getComputedStyle(ripple)['animation-duration'].replace('s', '')
  )

  const timeout = Math.max(200, duration * 1000 + delay * 1000)

  setTimeout(() => {
    ripple.classList.add('waves-ripple-out')

    setTimeout(() => {
      ripple.remove()
    }, 300)
  }, timeout)
}

function updateRipple(el, binding, wasEnabled) {
  const enabled = isRippleEnabled(binding.value)
  if (!enabled) {
    rippleHide({ currentTarget: el })
  }

  el[context] = el[context] || {}
  el[context].enabled = enabled

  const value = binding.value || {}
  if (value.center) {
    el[context].centered = true
  }
  if (value.class) {
    el[context].class = value.class
  }
  if (value.circle) {
    el[context].circle = value.circle
  }

  if (enabled && !wasEnabled) {
    el.addEventListener('mousedown', rippleShow, { passive: true })
    el.addEventListener('mouseup', rippleHide, { passive: true })
    el.addEventListener('mouseleave', rippleHide, { passive: true })
    // Dragstart is needed to prevent text selection in Safari
    el.addEventListener('dragstart', rippleHide, { passive: true })
  } else if (!enabled && wasEnabled) {
    removeListeners(el)
  }
}

function removeListeners(el) {
  el.removeEventListener('mousedown', rippleShow, { passive: true })
  el.removeEventListener('mouseup', rippleHide, { passive: true })
  el.removeEventListener('mouseleave', rippleHide, { passive: true })
  el.removeEventListener('dragstart', rippleHide, { passive: true })
}

function directive(el, binding) {
  updateRipple(el, binding, false)
}

function unbind(el) {
  delete el[context]
  removeListeners(el)
}

function update(el, binding) {
  if (binding.value === binding.oldValue) {
    return
  }

  const wasEnabled = isRippleEnabled(binding.oldValue)
  updateRipple(el, binding, wasEnabled)
}

export default {
  bind: directive,
  unbind,
  update
}