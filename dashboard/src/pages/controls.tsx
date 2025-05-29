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

    useEffect(() => {
        console.log("Coefficients data:", coefficients);
    }, [coefficients]);

    // Helper function to format coefficient names
    const formatCoefficientName = (key: string) => {
        return key
            .replace(/_coef$/, '') // Remove _coef suffix
            .replace(/_/g, ' ') // Replace underscores with spaces
            .replace(/\b\w/g, l => l.toUpperCase()); // Capitalize first letter of each word
    };

    // Get only the coefficient fields (ending with _coef)
    const getCoefficientFields = (data: any) => {
        return Object.entries(data)
            .filter(([key]) => key.endsWith('_coef'))
            .reduce((acc, [key, value]) => {
                acc[key] = value;
                return acc;
            }, {} as Record<string, any>);
    };

    return (
        <Card className="mb-4">
            <CardHeader>
                <CardTitle>Price Model Coefficients</CardTitle>
            </CardHeader>
            <CardContent>
                {loading && <div>Loading...</div>}
                {error && <div className="text-red-500">{error}</div>}
                {!loading && coefficients ? (
                    <div className="space-y-4">
                        {/* Model Info */}
                        <div className="bg-muted p-4 rounded-lg">
                            <div className="grid grid-cols-2 gap-4 text-sm">
                                <div>
                                    <span className="font-medium">Model Version:</span> {coefficients.model_version}
                                </div>
                                <div>
                                    <span className="font-medium">Training Date:</span> {new Date(coefficients.training_date).toLocaleDateString()}
                                </div>
                            </div>
                        </div>
                        
                        {/* Coefficients */}
                        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                            {Object.entries(getCoefficientFields(coefficients)).map(([key, value]) => (
                                <div key={key} className="p-3 bg-muted rounded-lg">
                                    <div className="text-sm font-medium text-muted-foreground">
                                        {formatCoefficientName(key)}
                                    </div>
                                    <div className="text-xl font-semibold">
                                        {typeof value === 'number' ? value.toFixed(6) : value}
                                    </div>
                                </div>
                            ))}
                        </div>
                    </div>
                ) : (
                    !loading && <div className="text-muted-foreground">No coefficients data available</div>
                )}
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
                <div className="grid grid-cols-2 gap-4">
                    <div className="space-y-4">
                        <div className="p-4 bg-muted rounded-lg">
                            <div className="text-lg font-semibold mb-1">R²</div>
                            <div className="text-2xl font-bold">{metrics.r_squared.toFixed(4)}</div>
                            <div className="text-sm text-muted-foreground mt-2">
                                Measures how well the model fits the data. Ranges from 0 to 1, where 1 indicates perfect prediction.
                            </div>
                        </div>
                        <div className="p-4 bg-muted rounded-lg">
                            <div className="text-lg font-semibold mb-1">RMSE</div>
                            <div className="text-2xl font-bold">{metrics.rmse.toFixed(4)}</div>
                            <div className="text-sm text-muted-foreground mt-2">
                                Root Mean Square Error - measures the average magnitude of prediction errors. Lower is better.
                            </div>
                        </div>
                    </div>
                    <div className="space-y-4">
                        <div className="p-4 bg-muted rounded-lg">
                            <div className="text-lg font-semibold mb-1">MAE</div>
                            <div className="text-2xl font-bold">{metrics.mae.toFixed(4)}</div>
                            <div className="text-sm text-muted-foreground mt-2">
                                Mean Absolute Error - measures the average absolute difference between predicted and actual values.
                            </div>
                        </div>
                        <div className="p-4 bg-muted rounded-lg">
                            <div className="text-lg font-semibold mb-1">Sample Size</div>
                            <div className="text-2xl font-bold">{metrics.sample_size.toLocaleString()}</div>
                            <div className="text-sm text-muted-foreground mt-2">
                                Number of data points used to train the model.
                            </div>
                        </div>
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
        <div className="space-y-4">
            <div className="flex gap-4">
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
            <div className="text-sm text-muted-foreground space-y-2">
                <div className="flex items-center gap-2">
                    <div className="w-2 h-2 bg-green-500 rounded-full"></div>
                    <span>Price adjustments run automatically daily at 00:00 UTC</span>
                </div>
                <div className="flex items-center gap-2">
                    <div className="w-2 h-2 bg-green-500 rounded-full"></div>
                    <span>Model training runs automatically on the 1st of each month at 01:00 UTC</span>
                </div>
            </div>
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