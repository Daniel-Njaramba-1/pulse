<script lang="ts">
    import { onMount } from 'svelte';

    // Demo price update simulation
    type Product = {
        id: string;
        name: string;
        basePrice: number;
        price: number;
        lastPrice: number;
        lastUpdated: Date;
        changeAmount: number;
        changePercent: number;
        isIncrease: boolean;
        lastUpdateSeconds: number;
    };

    let products: Product[] = [
        {
            id: 'macbook',
            name: 'MacBook Air M2',
            basePrice: 1199.99,
            price: 1199.99,
            lastPrice: 1199.99,
            lastUpdated: new Date(),
            changeAmount: -50,
            changePercent: -4.0,
            isIncrease: false,
            lastUpdateSeconds: 2
        },
        {
            id: 'iphone',
            name: 'iPhone 15 Pro',
            basePrice: 999.99,
            price: 999.99,
            lastPrice: 999.99,
            lastUpdated: new Date(),
            changeAmount: 25,
            changePercent: 2.6,
            isIncrease: true,
            lastUpdateSeconds: 5
        },
        {
            id: 'airpods',
            name: 'AirPods Pro',
            basePrice: 249.99,
            price: 249.99,
            lastPrice: 249.99,
            lastUpdated: new Date(),
            changeAmount: -15,
            changePercent: -5.7,
            isIncrease: false,
            lastUpdateSeconds: 1
        }
    ];

    let interval: ReturnType<typeof setInterval>;

    function updatePrices() {
        products = products.map((product) => {
            const change = (Math.random() - 0.5) * 0.1; // -5% to +5%
            const newPrice = +(product.basePrice * (1 + change)).toFixed(2);
            const changeAmount = +(newPrice - product.basePrice).toFixed(2);
            const changePercent = +(change * 100).toFixed(1);
            const isIncrease = changeAmount > 0;
            return {
                ...product,
                lastPrice: product.price,
                price: newPrice,
                lastUpdated: new Date(),
                changeAmount,
                changePercent,
                isIncrease,
                lastUpdateSeconds: 0
            };
        });
    }

    function updateSeconds() {
        products = products.map((product) => ({
            ...product,
            lastUpdateSeconds: Math.floor((Date.now() - product.lastUpdated.getTime()) / 1000)
        }));
    }

    onMount(() => {
        interval = setInterval(() => {
            updatePrices();
        }, 3000);

        const secondsInterval = setInterval(() => {
            updateSeconds();
        }, 1000);

        // Initial update
        setTimeout(updatePrices, 1000);

        return () => {
            clearInterval(interval);
            clearInterval(secondsInterval);
        };
    });
</script>

<style>
    * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }


        .container {
            max-width: 1200px;
            margin: 0 auto;
            padding: 0 24px;
        }

        /* ========================= HERO SECTION ========================= */
        .hero {
            background: linear-gradient(135deg, #0f172a 0%, #1e293b 50%, #334155 100%);
            color: white;
            padding: 80px 0 120px;
            position: relative;
            overflow: hidden;
        }

        .hero::before {
            content: '';
            position: absolute;
            top: 0;
            left: 0;
            right: 0;
            bottom: 0;
            background: url('data:image/svg+xml;utf8,<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100"><defs><pattern id="grid" width="10" height="10" patternUnits="userSpaceOnUse"><path d="M 10 0 L 0 0 0 10" fill="none" stroke="rgba(255,255,255,0.05)" stroke-width="0.5"/></pattern></defs><rect width="100" height="100" fill="url(%23grid)"/></svg>');
            opacity: 0.3;
        }

        .hero-content {
            position: relative;
            z-index: 2;
            text-align: center;
        }

        .hero h1 {
            font-size: 4rem;
            font-weight: 900;
            margin-bottom: 24px;
            background: linear-gradient(135deg, #ffffff 0%, #cbd5e1 100%);
            -webkit-background-clip: text;
            -webkit-text-fill-color: transparent;
            background-clip: text;
            letter-spacing: -0.02em;
        }

        .hero .subtitle {
            font-size: 1.5rem;
            margin-bottom: 32px;
            color: #cbd5e1;
            font-weight: 300;
        }

        .hero .description {
            font-size: 1.125rem;
            max-width: 600px;
            margin: 0 auto 48px;
            color: #e2e8f0;
            line-height: 1.7;
        }

        .pulse-indicator {
            display: inline-flex;
            align-items: center;
            gap: 12px;
            background: rgba(34, 197, 94, 0.1);
            border: 1px solid rgba(34, 197, 94, 0.3);
            padding: 12px 24px;
            border-radius: 50px;
            margin-bottom: 48px;
            animation: glow 2s ease-in-out infinite alternate;
        }

        .pulse-dot {
            width: 12px;
            height: 12px;
            background: #22c55e;
            border-radius: 50%;
            animation: pulse 2s infinite;
        }

        @keyframes pulse {
            0% { box-shadow: 0 0 0 0 rgba(34, 197, 94, 0.7); }
            70% { box-shadow: 0 0 0 10px rgba(34, 197, 94, 0); }
            100% { box-shadow: 0 0 0 0 rgba(34, 197, 94, 0); }
        }

        @keyframes glow {
            from { box-shadow: 0 0 20px rgba(34, 197, 94, 0.2); }
            to { box-shadow: 0 0 30px rgba(34, 197, 94, 0.4); }
        }

        /* ========================= FEATURES SECTION ========================= */
        .features {
            padding: 100px 0;
            background: white;
        }

        .section-header {
            text-align: center;
            margin-bottom: 80px;
        }

        .section-title {
            font-size: 3rem;
            font-weight: 800;
            margin-bottom: 24px;
            background: linear-gradient(135deg, #1e293b 0%, #334155 100%);
            -webkit-background-clip: text;
            -webkit-text-fill-color: transparent;
            background-clip: text;
        }

        .section-subtitle {
            font-size: 1.25rem;
            color: #64748b;
            max-width: 600px;
            margin: 0 auto;
        }

        .features-grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(350px, 1fr));
            gap: 40px;
            margin-bottom: 80px;
        }

        .feature-card {
            background: linear-gradient(135deg, #ffffff 0%, #f8fafc 100%);
            border: 1px solid #e2e8f0;
            border-radius: 20px;
            padding: 40px;
            position: relative;
            overflow: hidden;
            transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
        }

        .feature-card:hover {
            transform: translateY(-8px);
            box-shadow: 0 20px 40px rgba(0, 0, 0, 0.1);
            border-color: #cbd5e1;
        }

        .feature-icon {
            width: 60px;
            height: 60px;
            background: linear-gradient(135deg, #3b82f6 0%, #1d4ed8 100%);
            border-radius: 16px;
            display: flex;
            align-items: center;
            justify-content: center;
            margin-bottom: 24px;
            font-size: 24px;
            color: white;
        }

        .feature-card h3 {
            font-size: 1.5rem;
            font-weight: 700;
            margin-bottom: 16px;
            color: #1e293b;
        }

        .feature-card p {
            color: #64748b;
            line-height: 1.7;
            margin-bottom: 20px;
        }

        /* ========================= HOW IT WORKS SECTION ========================= */
        .how-it-works {
            padding: 100px 0;
            background: linear-gradient(135deg, #f1f5f9 0%, #e2e8f0 100%);
        }

        .steps-container {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
            gap: 60px;
            margin-top: 60px;
        }

        .step {
            text-align: center;
            position: relative;
        }

        .step-number {
            width: 80px;
            height: 80px;
            background: linear-gradient(135deg, #6366f1 0%, #4338ca 100%);
            border-radius: 50%;
            display: flex;
            align-items: center;
            justify-content: center;
            font-size: 2rem;
            font-weight: 800;
            color: white;
            margin: 0 auto 24px;
            box-shadow: 0 10px 30px rgba(99, 102, 241, 0.3);
        }

        .step h3 {
            font-size: 1.5rem;
            font-weight: 700;
            margin-bottom: 16px;
            color: #1e293b;
        }

        .step p {
            color: #64748b;
            line-height: 1.7;
        }

        /* ========================= DEMO SECTION ========================= */
        .demo {
            padding: 100px 0;
            background: white;
        }

        .demo-container {
            background: linear-gradient(135deg, #1e293b 0%, #334155 100%);
            border-radius: 24px;
            padding: 60px;
            color: white;
            position: relative;
            overflow: hidden;
        }

        .demo-header {
            text-align: center;
            margin-bottom: 48px;
        }

        .demo-title {
            font-size: 2.5rem;
            font-weight: 800;
            margin-bottom: 16px;
        }

        .price-demo {
            display: flex;
            justify-content: center;
            align-items: center;
            gap: 40px;
            flex-wrap: wrap;
        }

        .price-card {
            background: rgba(255, 255, 255, 0.1);
            backdrop-filter: blur(10px);
            border: 1px solid rgba(255, 255, 255, 0.2);
            border-radius: 16px;
            padding: 32px;
            text-align: center;
            min-width: 250px;
            transition: all 0.3s ease;
        }

        .price-card:hover {
            transform: scale(1.05);
            background: rgba(255, 255, 255, 0.15);
        }

        .product-name {
            font-size: 1.25rem;
            font-weight: 600;
            margin-bottom: 16px;
        }

        .price-display {
            font-size: 3rem;
            font-weight: 900;
            color: #22c55e;
            margin-bottom: 8px;
            transition: all 0.5s ease;
        }

        .price-change {
            font-size: 0.875rem;
            padding: 4px 12px;
            border-radius: 20px;
            font-weight: 600;
        }

        .price-up {
            background: rgba(239, 68, 68, 0.2);
            color: #fca5a5;
        }

        .price-down {
            background: rgba(34, 197, 94, 0.2);
            color: #86efac;
        }

        .last-updated {
            margin-top: 16px;
            font-size: 0.75rem;
            color: #cbd5e1;
        }

        /* ========================= CTA SECTION ========================= */
        .cta {
            padding: 100px 0;
            background: linear-gradient(135deg, #0f172a 0%, #1e293b 100%);
            color: white;
            text-align: center;
        }

        .cta-button {
            display: inline-flex;
            align-items: center;
            gap: 12px;
            background: linear-gradient(135deg, #3b82f6 0%, #1d4ed8 100%);
            color: white;
            padding: 20px 40px;
            border-radius: 50px;
            text-decoration: none;
            font-weight: 700;
            font-size: 1.125rem;
            transition: all 0.3s ease;
            box-shadow: 0 10px 30px rgba(59, 130, 246, 0.3);
        }

        .cta-button:hover {
            transform: translateY(-2px);
            box-shadow: 0 15px 40px rgba(59, 130, 246, 0.4);
            background: linear-gradient(135deg, #2563eb 0%, #1d4ed8 100%);
        }

        /* ========================= RESPONSIVE DESIGN ========================= */
        @media (max-width: 768px) {
            .hero h1 {
                font-size: 2.5rem;
            }
            
            .hero .subtitle {
                font-size: 1.25rem;
            }
            
            .section-title {
                font-size: 2rem;
            }
            
            .features-grid {
                grid-template-columns: 1fr;
                gap: 30px;
            }
            
            .feature-card {
                padding: 30px;
            }
            
            .steps-container {
                grid-template-columns: 1fr;
                gap: 40px;
            }
            
            .price-demo {
                flex-direction: column;
                gap: 20px;
            }
            
            .demo-container {
                padding: 40px 20px;
            }
        }
</style>

<!-- ========================= HERO SECTION ========================= -->
<section class="hero">
    <div class="container">
        <div class="hero-content">
            <div class="pulse-indicator">
                <div class="pulse-dot"></div>
                <span>Live Price Updates Active</span>
            </div>
            <h1>
                <span style="font-weight:900;">Dynamic Pricing</span><br>
                <span style="font-weight:900;">Real-Time Transparency</span>
            </h1>
            <p class="subtitle">Experience the future of fair pricing</p>
            <p class="description">
                Our revolutionary <b>Server-Sent Events (SSE)</b> technology delivers instant price updates,
                ensuring you always see the most current market rates with complete transparency.
            </p>
        </div>
    </div>
</section>

<!-- ========================= FEATURES SECTION ========================= -->
<section class="features">
    <div class="container">
        <div class="section-header">
            <h2 class="section-title">Why Dynamic Pricing?</h2>
            <p class="section-subtitle">
                Get the best deals with real-time market adjustments and complete pricing transparency
            </p>
        </div>
        <div class="features-grid">
            <div class="feature-card">
                <div class="feature-icon">⚡</div>
                <h3>Instant Updates</h3>
                <p>
                    Prices update in real-time using <b>Server-Sent Events (SSE)</b> technology.
                    No page refresh needed - watch prices adjust live as market conditions change.
                </p>
            </div>
            <div class="feature-card">
                <div class="feature-icon">🎯</div>
                <h3>Fair Market Pricing</h3>
                <p>
                    Our dynamic system ensures competitive prices that reflect real market value,
                    protecting you from overpricing while rewarding early adopters.
                </p>
            </div>
        </div>
    </div>
</section>

<!-- ========================= HOW IT WORKS SECTION ========================= -->
<section class="how-it-works">
    <div class="container">
        <div class="section-header">
            <h2 class="section-title">How It Works</h2>
            <p class="section-subtitle">
                Our advanced SSE technology delivers seamless real-time pricing updates
            </p>
        </div>
        <div class="steps-container">
            <div class="step">
                <div class="step-number">1</div>
                <h3>Connect & Monitor</h3>
                <p>
                    When you visit our site, we establish a <b>Server-Sent Event</b> connection
                    that monitors real-time price changes for all products you're viewing.
                </p>
            </div>
            <div class="step">
                <div class="step-number">2</div>
                <h3>Real-Time Analysis</h3>
                <p>
                    Our system continuously analyzes market demand, inventory levels,
                    competitor pricing, and seasonal trends to calculate optimal prices.
                </p>
            </div>
            <div class="step">
                <div class="step-number">3</div>
                <h3>Instant Updates</h3>
                <p>
                    Price changes are pushed directly to your browser via <b>SSE</b>,
                    showing you live updates with smooth animations and change indicators.
                </p>
            </div>
            <div class="step">
                <div class="step-number">4</div>
                <h3>Make Informed Decisions</h3>
                <p>
                    Use our price history charts, trend indicators, and transparency
                    data to make confident purchasing decisions at the right time.
                </p>
            </div>
        </div>
    </div>
</section>

<!-- ========================= DEMO SECTION ========================= -->
<section class="demo">
    <div class="container">
        <div class="demo-container">
            <div class="demo-header">
                <h2 class="demo-title">Live Pricing Demo</h2>
                <p>Watch real-time price updates in action</p>
            </div>
            <div class="price-demo">
                {#each products as product}
                    <div class="price-card">
                        <div class="product-name">{product.name}</div>
                        <div
                            class="price-display"
                            style="transition: all 0.5s; color: #22c55e;"
                        >
                            ${product.price.toFixed(2)}
                        </div>
                        <div class="price-change {product.isIncrease ? 'price-up' : 'price-down'}">
                            {product.isIncrease ? '↑' : '↓'} ${Math.abs(product.changeAmount).toFixed(0)} ({Math.abs(product.changePercent)}%)
                        </div>
                        <div class="last-updated">
                            Updated {product.lastUpdateSeconds} second{product.lastUpdateSeconds === 1 ? '' : 's'} ago
                        </div>
                    </div>
                {/each}
            </div>
        </div>
    </div>
</section>

<!-- ========================= CTA SECTION ========================= -->
<section class="cta">
    <div class="container">
        <h2 class="section-title">Ready to Experience Dynamic Pricing?</h2>
        <p class="section-subtitle" style="color: #cbd5e1; margin-bottom: 40px;">
            Join thousands of smart shoppers getting the best deals with real-time transparency
        </p>
        <a href="/" class="cta-button">
            <span>Start Shopping Now</span>
            <span>→</span>
        </a>
    </div>
</section>

