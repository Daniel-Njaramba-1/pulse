import { useEffect } from "react";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { ControlsProvider, useControls } from "@/contexts/controls";

function CoefficientsCard() {
    const { coefficients, loading, error, fetchCoefficients } = useControls();

    useEffect(() => {
        fetchCoefficients();
        // eslint-disable-next-line
    }, []);

    return (
        <Card className="mb-4">
            <CardHeader>
                <CardTitle>Price Model Coefficients</CardTitle>
            </CardHeader>
            <CardContent>
                {loading && <div>Loading...</div>}
                {error && <div className="text-red-500">{error}</div>}
                {!loading && coefficients && coefficients.length > 0 ? (
                    <div className="space-y-2">
                        {coefficients.map((coef, idx) => (
                            <div key={idx} className="flex flex-wrap gap-4">
                                {Object.entries(coef).map(([key, value]) => (
                                    <div key={key} className="text-sm">
                                        <span className="font-medium">{key}:</span> {value}
                                    </div>
                                ))}
                            </div>
                        ))}
                    </div>
                ) : null}
            </CardContent>
        </Card>
    );
}

function MetricsCard() {
    const { metrics } = useControls();

    if (!metrics) return null;

    return (
        <Card className="mb-4">
            <CardHeader>
                <CardTitle>Model Metrics</CardTitle>
            </CardHeader>
            <CardContent>
                <div className="space-y-1 text-sm">
                    <div>
                        <span className="font-medium">R²:</span> {metrics.r_squared}
                    </div>
                    <div>
                        <span className="font-medium">RMSE:</span> {metrics.rmse}
                    </div>
                    <div>
                        <span className="font-medium">MAE:</span> {metrics.mae}
                    </div>
                    <div>
                        <span className="font-medium">Sample Size:</span> {metrics.sample_size}
                    </div>
                </div>
            </CardContent>
        </Card>
    );
}

function ControlsButtons() {
    const {
        loading,
        triggerAdjustAllPrices,
        triggerTrainModel,
        fetchCoefficients,
    } = useControls();

    return (
        <div className="flex gap-4 mb-4">
            <Button onClick={triggerAdjustAllPrices} disabled={loading}>
                Adjust All Prices
            </Button>
            <Button onClick={async () => {
                await triggerTrainModel();
                await fetchCoefficients();
            }} disabled={loading}>
                Train Model
            </Button>
        </div>
    );
}

export default function ControlsPage() {
    return (
        <ControlsProvider>
            <div className="max-w-2xl mx-auto py-8 px-4">
                <ControlsButtons />
                <CoefficientsCard />
                <MetricsCard />
            </div>
        </ControlsProvider>
    );
}