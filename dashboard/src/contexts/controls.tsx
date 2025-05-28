import { createContext, useContext, useState, type ReactNode } from "react";

type Coefficient = {
    [key: string]: number | string;
};

type ModelMetrics = {
    r_squared: number;
    rmse: number;
    mae: number;
    sample_size: number;
};

type ControlsContextType = {
    coefficients: Coefficient[] | null;
    metrics: ModelMetrics | null;
    loading: boolean;
    error: string | null;
    fetchCoefficients: () => Promise<void>;
    triggerAdjustAllPrices: () => Promise<void>;
    triggerTrainModel: () => Promise<void>;
};

const ControlsContext = createContext<ControlsContextType | undefined>(undefined);

export const ControlsProvider = ({ children }: { children: ReactNode }) => {
    const [coefficients, setCoefficients] = useState<Coefficient[] | null>(null);
    const [metrics, setMetrics] = useState<ModelMetrics | null>(null);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    const fetchCoefficients = async () => {
        setLoading(true);
        setError(null);
        try {
            const res = await fetch("http://localhost:5872/get-price-model-coefficients");
            const data = await res.json();
            if (data.coefficients) setCoefficients(data.coefficients);
            else setError("No coefficients found");
        } catch (e: any) {
            setError(e.message || "Failed to fetch coefficients");
        } finally {
            setLoading(false);
        }
    };

    const triggerAdjustAllPrices = async () => {
        setLoading(true);
        setError(null);
        try {
            const res = await fetch("http://localhost:5872/adjust-prices", { method: "POST" });
            if (!res.ok) throw new Error("Failed to adjust all prices");
        } catch (e: any) {
            setError(e.message || "Failed to adjust all prices");
        } finally {
            setLoading(false);
        }
    };

    const triggerTrainModel = async () => {
        setLoading(true);
        setError(null);
        try {
            const res = await fetch("http://localhost:5872/train-model", { method: "POST" });
            const data = await res.json();
            if (data.metrics) setMetrics(data.metrics);
            else setError(data.message || "Failed to train model");
        } catch (e: any) {
            setError(e.message || "Failed to train model");
        } finally {
            setLoading(false);
        }
    };

    return (
        <ControlsContext.Provider
            value={{
                coefficients,
                metrics,
                loading,
                error,
                fetchCoefficients,
                triggerAdjustAllPrices,
                triggerTrainModel,
            }}
        >
            {children}
        </ControlsContext.Provider>
    );
};

export const useControls = () => {
    const ctx = useContext(ControlsContext);
    if (!ctx) throw new Error("useControls must be used within ControlsProvider");
    return ctx;
};