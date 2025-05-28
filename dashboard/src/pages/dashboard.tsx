import { useEffect } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import {
    LineChart, Line, BarChart, Bar, PieChart, Pie, Cell, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer
} from 'recharts';
import { useDashboard } from '@/contexts/dashboard';

const COLORS = ['#8884d8', '#82ca9d', '#ffc658', '#ff8042', '#0088FE', '#00C49F', '#FFBB28', '#FF8042'];

export default function Dashboard() {
    const {
        dashboardData,
        loading,
        error,
        initializeDashboard,
    } = useDashboard();

    useEffect(() => {
        initializeDashboard();
    }, [initializeDashboard]);

    if (loading) return <div>Loading...</div>;
    if (error) return <div className="text-red-500">{error}</div>;
    if (!dashboardData) return null;

    return (
        <div className="space-y-6">
            <Tabs defaultValue="sales">
                <TabsList>
                    <TabsTrigger value="sales">Sales</TabsTrigger>
                    <TabsTrigger value="inventory">Inventory</TabsTrigger>
                    <TabsTrigger value="pricing">Pricing</TabsTrigger>
                    <TabsTrigger value="customers">Customers</TabsTrigger>
                    <TabsTrigger value="categories">Categories</TabsTrigger>
                </TabsList>
                <TabsContent value="sales">
                    <Card>
                        <CardHeader>
                            <CardTitle>Sales Analytics</CardTitle>
                        </CardHeader>
                        <CardContent>
                            <ResponsiveContainer width="100%" height={300}>
                                <LineChart data={dashboardData.sales_analytics}>
                                    <CartesianGrid strokeDasharray="3 3" />
                                    <XAxis dataKey="date" />
                                    <YAxis />
                                    <Tooltip />
                                    <Legend />
                                    <Line type="monotone" dataKey="total_sales" stroke="#8884d8" name="Total Sales" />
                                    <Line type="monotone" dataKey="total_revenue" stroke="#82ca9d" name="Total Revenue" />
                                </LineChart>
                            </ResponsiveContainer>
                        </CardContent>
                    </Card>
                </TabsContent>
                <TabsContent value="inventory">
                    <Card>
                        <CardHeader>
                            <CardTitle>Inventory Status</CardTitle>
                        </CardHeader>
                        <CardContent>
                            <ResponsiveContainer width="100%" height={300}>
                                <BarChart data={dashboardData.inventory_status}>
                                    <CartesianGrid strokeDasharray="3 3" />
                                    <XAxis dataKey="product_name" />
                                    <YAxis />
                                    <Tooltip />
                                    <Legend />
                                    <Bar dataKey="current_stock" fill="#8884d8" name="Current Stock" />
                                    <Bar dataKey="stock_threshold" fill="#ffc658" name="Stock Threshold" />
                                </BarChart>
                            </ResponsiveContainer>
                        </CardContent>
                    </Card>
                </TabsContent>
                <TabsContent value="pricing">
                    <Card>
                        <CardHeader>
                            <CardTitle>Pricing Analytics</CardTitle>
                        </CardHeader>
                        <CardContent>
                            <ResponsiveContainer width="100%" height={300}>
                                <BarChart data={dashboardData.pricing_analytics}>
                                    <CartesianGrid strokeDasharray="3 3" />
                                    <XAxis dataKey="product_name" />
                                    <YAxis />
                                    <Tooltip />
                                    <Legend />
                                    <Bar dataKey="base_price" fill="#8884d8" name="Base Price" />
                                    <Bar dataKey="adjusted_price" fill="#82ca9d" name="Adjusted Price" />
                                </BarChart>
                            </ResponsiveContainer>
                        </CardContent>
                    </Card>
                </TabsContent>
                <TabsContent value="customers">
                    <Card>
                        <CardHeader>
                            <CardTitle>Customer Behavior</CardTitle>
                        </CardHeader>
                        <CardContent>
                            <ResponsiveContainer width="100%" height={300}>
                                <BarChart data={dashboardData.customer_behavior}>
                                    <CartesianGrid strokeDasharray="3 3" />
                                    <XAxis dataKey="product_name" />
                                    <YAxis />
                                    <Tooltip />
                                    <Legend />
                                    <Bar dataKey="average_rating" fill="#8884d8" name="Avg Rating" />
                                    <Bar dataKey="review_count" fill="#ffc658" name="Review Count" />
                                    <Bar dataKey="wishlist_count" fill="#82ca9d" name="Wishlist Count" />
                                </BarChart>
                            </ResponsiveContainer>
                        </CardContent>
                    </Card>
                </TabsContent>
                <TabsContent value="categories">
                    <Card>
                        <CardHeader>
                            <CardTitle>Category Revenue</CardTitle>
                        </CardHeader>
                        <CardContent>
                            <ResponsiveContainer width="100%" height={300}>
                                <PieChart>
                                    <Pie
                                        data={dashboardData.category_revenue}
                                        dataKey="total_revenue"
                                        nameKey="category_name"
                                        cx="50%"
                                        cy="50%"
                                        outerRadius={100}
                                        label
                                    >
                                        {dashboardData.category_revenue.map((_, index) => (
                                            <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                                        ))}
                                    </Pie>
                                    <Tooltip />
                                    <Legend />
                                </PieChart>
                            </ResponsiveContainer>
                        </CardContent>
                    </Card>
                </TabsContent>
            </Tabs>
        </div>
    );
}