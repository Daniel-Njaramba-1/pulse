<script lang="ts">
  import { goto } from '$app/navigation';
    import { onMount } from 'svelte';

    // Smooth scrolling for anchor links
    onMount(() => {
        const handleAnchorClick = (e: Event) => {
            const href = (e.currentTarget as HTMLAnchorElement).getAttribute('href');
            if (href && href.startsWith('#')) {
                e.preventDefault();
                const target = document.querySelector(href);
                if (target) {
                    target.scrollIntoView({
                        behavior: 'smooth',
                        block: 'start'
                    });
                }
            }
        };

        document.querySelectorAll('a[href^="#"]').forEach(anchor => {
            anchor.addEventListener('click', handleAnchorClick);
        });

        // Intersection observer for feature cards
        const observerOptions = {
            threshold: 0.1,
            rootMargin: '0px 0px -50px 0px'
        };

        const observer = new IntersectionObserver((entries) => {
            entries.forEach(entry => {
                if (entry.isIntersecting) {
                    (entry.target as HTMLElement).style.animationPlayState = 'running';
                }
            });
        }, observerOptions);

        document.querySelectorAll('.feature-card').forEach(card => {
            observer.observe(card);
        });

        // Dynamic background particles (simplified)
        function createParticle() {
            const particle = document.createElement('div');
            particle.className = 'shape';
            const size = Math.random() * 50 + 30;
            particle.style.width = size + 'px';
            particle.style.height = size + 'px';
            particle.style.left = Math.random() * 100 + '%';
            particle.style.top = '100%';
            particle.style.animationDuration = (Math.random() * 10 + 15) + 's';

            const shapesContainer = document.querySelector('.floating-shapes');
            if (shapesContainer) {
                shapesContainer.appendChild(particle);
            }

            setTimeout(() => {
                particle.remove();
            }, 25000);
        }

        const interval = setInterval(createParticle, 8000);

        return () => {
            clearInterval(interval);
            document.querySelectorAll('a[href^="#"]').forEach(anchor => {
                anchor.removeEventListener('click', handleAnchorClick);
            });
        };
    });
</script>

<style>
    :global(html, body) {
        height: 100%;
        margin: 0;
        padding: 0;
    }

    * {
        margin: 0;
        padding: 0;
        box-sizing: border-box;
    }

    /* body {
        font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
        background: #f5f6fa;
        min-height: 100vh;
        overflow-x: hidden;
    } */

    .animated-bg {
        position: fixed;
        top: 0;
        left: 0;
        width: 100%;
        height: 100%;
        z-index: -1;
        background: linear-gradient(135deg, #f5f6fa 0%, #e9ecef 100%);
        background-size: 400% 400%;
        animation: gradientShift 15s ease infinite;
    }

    @keyframes gradientShift {
        0% { background-position: 0% 50%; }
        50% { background-position: 100% 50%; }
        100% { background-position: 0% 50%; }
    }

    .floating-shapes {
        position: fixed;
        top: 0;
        left: 0;
        width: 100%;
        height: 100%;
        z-index: -1;
        pointer-events: none;
    }

    .shape {
        position: absolute;
        background: rgba(120, 120, 120, 0.08);
        border-radius: 50%;
        animation: float 20s infinite linear;
    }

    .shape:nth-child(1) {
        width: 80px;
        height: 80px;
        top: 20%;
        left: 10%;
        animation-delay: 0s;
    }

    .shape:nth-child(2) {
        width: 120px;
        height: 120px;
        top: 60%;
        right: 15%;
        animation-delay: -5s;
    }

    .shape:nth-child(3) {
        width: 60px;
        height: 60px;
        top: 80%;
        left: 70%;
        animation-delay: -10s;
    }

    @keyframes float {
        0% { transform: translateY(0px) rotate(0deg); }
        33% { transform: translateY(-30px) rotate(120deg); }
        66% { transform: translateY(30px) rotate(240deg); }
        100% { transform: translateY(0px) rotate(360deg); }
    }

    .container {
        max-width: 1200px;
        margin: 0 auto;
        padding: 2rem;
        position: relative;
        z-index: 1;
    }

    .hero-section {
        text-align: center;
        margin-bottom: 4rem;
        padding: 4rem 2rem;
    }

    .hero-title {
        font-size: clamp(2.5rem, 6vw, 4rem);
        font-weight: 800;
        background: linear-gradient(135deg, #222 0%, #666 100%);
        -webkit-background-clip: text;
        -webkit-text-fill-color: transparent;
        background-clip: text;
        margin-bottom: 1.5rem;
        text-shadow: 0 4px 8px rgba(0,0,0,0.08);
        animation: fadeInUp 1s ease-out;
    }

    .hero-subtitle {
        font-size: 1.25rem;
        color: #444;
        margin-bottom: 3rem;
        max-width: 600px;
        margin-left: auto;
        margin-right: auto;
        line-height: 1.7;
        animation: fadeInUp 1s ease-out 0.2s both;
    }

    .cta-buttons {
        display: flex;
        gap: 1.5rem;
        justify-content: center;
        flex-wrap: wrap;
        animation: fadeInUp 1s ease-out 0.4s both;
    }

    .btn {
        padding: 1rem 2.5rem;
        border: none;
        border-radius: 50px;
        font-size: 1.1rem;
        font-weight: 600;
        cursor: pointer;
        transition: all 0.3s ease;
        text-decoration: none;
        display: inline-flex;
        align-items: center;
        gap: 0.5rem;
        position: relative;
        overflow: hidden;
    }

    .btn-primary {
        background: linear-gradient(135deg, #ff6b6b, #ff8e8e);
        color: white;
        box-shadow: 0 8px 25px rgba(255, 107, 107, 0.2);
    }

    .btn-secondary {
        background: rgba(255, 255, 255, 0.7);
        color: #222;
        border: 2px solid rgba(200, 200, 200, 0.3);
        backdrop-filter: blur(10px);
    }

    .btn:hover {
        transform: translateY(-3px);
        box-shadow: 0 12px 35px rgba(0,0,0,0.08);
    }

    .btn-primary:hover {
        box-shadow: 0 12px 35px rgba(255, 107, 107, 0.3);
    }

    .features-grid {
        display: grid;
        grid-template-columns: repeat(auto-fit, minmax(350px, 1fr));
        gap: 2rem;
        margin-bottom: 4rem;
    }

    .feature-card {
        background: rgba(255, 255, 255, 0.7);
        backdrop-filter: blur(20px);
        border-radius: 20px;
        padding: 2.5rem;
        border: 1px solid rgba(200, 200, 200, 0.2);
        transition: all 0.4s ease;
        position: relative;
        overflow: hidden;
        animation: fadeInUp 1s ease-out;
        animation-play-state: paused;
    }

    .feature-card::before {
        content: '';
        position: absolute;
        top: 0;
        left: -100%;
        width: 100%;
        height: 100%;
        background: linear-gradient(90deg, transparent, rgba(255,255,255,0.1), transparent);
        transition: left 0.6s ease;
    }

    .feature-card:hover::before {
        left: 100%;
    }

    .feature-card:hover {
        transform: translateY(-10px) scale(1.02);
        box-shadow: 0 20px 40px rgba(0,0,0,0.08);
        border-color: rgba(200, 200, 200, 0.4);
    }

    .feature-icon {
        width: 70px;
        height: 70px;
        background: linear-gradient(135deg, #ff6b6b, #ff8e8e);
        border-radius: 16px;
        display: flex;
        align-items: center;
        justify-content: center;
        font-size: 2rem;
        margin-bottom: 1.5rem;
        box-shadow: 0 8px 25px rgba(255, 107, 107, 0.1);
    }

    .feature-title {
        font-size: 1.5rem;
        font-weight: 700;
        color: #222;
        margin-bottom: 1rem;
    }

    .feature-description {
        color: #444;
        line-height: 1.6;
        margin-bottom: 1.5rem;
    }

    .feature-list {
        list-style: none;
        padding: 0;
    }

    .feature-list li {
        color: #333;
        padding: 0.5rem 0;
        position: relative;
        padding-left: 1.5rem;
    }

    .feature-list li::before {
        content: '✦';
        position: absolute;
        left: 0;
        color: #ff6b6b;
        font-weight: bold;
    }

    /* .stats-section {
        background: rgba(255, 255, 255, 0.7);
        backdrop-filter: blur(20px);
        border-radius: 20px;
        padding: 3rem;
        margin-bottom: 4rem;
        border: 1px solid rgba(200, 200, 200, 0.1);
        animation: fadeInUp 1s ease-out 0.6s both;
    }

    .stats-grid {
        display: grid;
        grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
        gap: 2rem;
        text-align: center;
    }

    .stat-item {
        color: #222;
    }

    .stat-number {
        font-size: 2.5rem;
        font-weight: 800;
        background: linear-gradient(135deg, #ff6b6b, #ff8e8e);
        -webkit-background-clip: text;
        -webkit-text-fill-color: transparent;
        background-clip: text;
        display: block;
    }

    .stat-label {
        font-size: 1rem;
        opacity: 0.8;
        margin-top: 0.5rem;
    } */

    @keyframes fadeInUp {
        from {
            opacity: 0;
            transform: translateY(30px);
        }
        to {
            opacity: 1;
            transform: translateY(0);
        }
    }

    @media (max-width: 768px) {
        .cta-buttons {
            flex-direction: column;
            align-items: center;
        }

        .btn {
            width: 100%;
            max-width: 300px;
            justify-content: center;
        }

        .features-grid {
            grid-template-columns: 1fr;
        }

        /* .stats-grid {
            grid-template-columns: repeat(2, 1fr);
        } */
    }

    .pulse {
        animation: pulse 2s infinite;
    }

    @keyframes pulse {
        0% { transform: scale(1); }
        50% { transform: scale(1.05); }
        100% { transform: scale(1); }
    }
</style>

<div class="animated-bg"></div>

<div class="floating-shapes">
    <div class="shape"></div>
    <div class="shape"></div>
    <div class="shape"></div>
</div>

<div class="container">
    <section class="hero-section">
        <h1 class="hero-title">Admin Control Center</h1>
        <p class="hero-subtitle">
            Unleash the power of modern administration with our cutting-edge dashboard.
            Manage brands, categories, and products with unprecedented ease and efficiency.
        </p>
        <div class="cta-buttons">
            <a href="http://localhost:5185/" class="btn btn-primary pulse" target="_blank" rel="noopener">
                📊 View Graph Reports
            </a>
            <a href="#features" class="btn btn-secondary">
                🚀 Explore Features
            </a>
        </div>
    </section>

    <section id="features" class="features-grid">
        <button
            type="button"
            class="feature-card"
            aria-label="Go to Brand Management"
            onclick={() => goto(`/brands`)}
            style="cursor: pointer; background: none; border: none; padding: 0; text-align: left;"
        >
            <div class="feature-icon">🏢</div>
            <h3 class="feature-title">Brand Management</h3>
            <p class="feature-description">
            Transform your brand portfolio with intelligent management tools designed for the modern enterprise.
            </p>
            <ul class="feature-list">
            <li>Create stunning brand profiles with rich media support</li>
            <li>Real-time editing with instant preview capabilities</li>
            <li>Advanced search and filtering with AI-powered suggestions</li>
            <li>Brand performance analytics and insights</li>
            </ul>
        </button>

        <button
            type="button"
            class="feature-card"
            aria-label="Go to Category Management"
            onclick={() => goto(`/cateogires`)}
            style="cursor: pointer; background: none; border: none; padding: 0; text-align: left;"
        >
            <div class="feature-icon">📂</div>
            <h3 class="feature-title">Category Management</h3>
            <p class="feature-description">
                Organize your inventory with precision using our intelligent categorization system.
            </p>
            <ul class="feature-list">
                <li>Drag-and-drop category hierarchy builder</li>
                <li>Smart categorization with ML-powered suggestions</li>
                <li>Visual tree view with infinite nesting levels</li>
                <li>Bulk operations and mass category updates</li>
            </ul>
        </button>

        <button
            type="button"
            class="feature-card"
            aria-label="Go to Product Management"
            onclick={() => goto(`/brands`)}
            style="cursor: pointer; background: none; border: none; padding: 0; text-align: left;"
        >
            <div class="feature-icon">📦</div>
            <h3 class="feature-title">Product Management</h3>
            <p class="feature-description">
                Elevate your product catalog with comprehensive management tools that scale with your business.
            </p>
            <ul class="feature-list">
                <li>Rich product editor with image galleries</li>
                <li>Real-time inventory tracking and alerts</li>
                <li>Advanced pricing strategies and bulk updates</li>
                <li>Product performance metrics and optimization</li>
            </ul>
        </button>
    </section>
</div>
