`timescale 1ns / 1ps

// =====================================================================
// SOMA OS: Top-Level Virtualizer (Clean Edition V4)
// =====================================================================

module top_quantum_virtualizer (
    // Zynq Fixed IO (Required for PS stability - Matched to ALINX 7020)
    inout wire [14:0] DDR_addr,
    inout wire [2:0] DDR_ba,
    inout wire DDR_cas_n,
    inout wire DDR_ck_n,
    inout wire DDR_ck_p,
    inout wire DDR_cke,
    inout wire DDR_cs_n,
    inout wire [3:0] DDR_dm,
    inout wire [31:0] DDR_dq,
    inout wire [3:0] DDR_dqs_n,
    inout wire [3:0] DDR_dqs_p,
    inout wire DDR_odt,
    inout wire DDR_ras_n,
    inout wire DDR_reset_n,
    inout wire DDR_we_n,
    inout wire FIXED_IO_ddr_vrn,
    inout wire FIXED_IO_ddr_vrp,
    inout wire [53:0] FIXED_IO_mio,
    inout wire FIXED_IO_ps_clk,
    inout wire FIXED_IO_ps_porb,
    inout wire FIXED_IO_ps_srstb,

    // Physical Indicators
    output wire led_0, 
    output wire led_1,
    output wire led_2,
    output wire led_3,
    
    // FPGA Internal Thermal Sensors (XADC)
    input wire vauxp0,          
    input wire vauxn0,
    input wire CLK100MHZ // External clock pin H16
);

    wire [15:0] xadc_temp_data_full;  
    wire [11:0] xadc_temp_data;       
    wire [11:0] calculated_psi_sc;    
    wire trigger_i2c;                 
    wire phase_field_active;
    wire master_entanglement_bus;
    wire sys_clk;

    // 0. The Zynq-7000 Processing System Wrapper
    design_1_wrapper zynq_bd (
        .DDR_0_addr(DDR_addr),
        .DDR_0_ba(DDR_ba),
        .DDR_0_cas_n(DDR_cas_n),
        .DDR_0_ck_n(DDR_ck_n),
        .DDR_0_ck_p(DDR_ck_p),
        .DDR_0_cke(DDR_cke),
        .DDR_0_cs_n(DDR_cs_n),
        .DDR_0_dm(DDR_dm),
        .DDR_0_dq(DDR_dq),
        .DDR_0_dqs_n(DDR_dqs_n),
        .DDR_0_dqs_p(DDR_dqs_p),
        .DDR_0_odt(DDR_odt),
        .DDR_0_ras_n(DDR_ras_n),
        .DDR_0_reset_n(DDR_reset_n),
        .DDR_0_we_n(DDR_we_n),
        .FIXED_IO_0_ddr_vrn(FIXED_IO_ddr_vrn),
        .FIXED_IO_0_ddr_vrp(FIXED_IO_ddr_vrp),
        .FIXED_IO_0_mio(FIXED_IO_mio),
        .FIXED_IO_0_ps_clk(FIXED_IO_ps_clk),
        .FIXED_IO_0_ps_porb(FIXED_IO_ps_porb),
        .FIXED_IO_0_ps_srstb(FIXED_IO_ps_srstb),
        .FCLK_CLK0(sys_clk)
    );

    // 1. Internal Observer: Digilent XADC Core
    xadc_wiz_0 internal_observer (
        .daddr_in(7'h00),             
        .den_in(1'b1),                
        .di_in(16'b0), 
        .dwe_in(1'b0), 
        .do_out(xadc_temp_data_full),      
        .drdy_out(trigger_i2c),       
        .dclk_in(sys_clk), 
        .vp_in(1'b0), 
        .vn_in(1'b0),
        .vauxp0(vauxp0),              
        .vauxn0(vauxn0)
    );

    assign xadc_temp_data = xadc_temp_data_full[15:4];

    // 2. Phase Field Activation Logic
    assign phase_field_active = (xadc_temp_data > 12'h800) ? 1'b1 : 1'b0;

    // 3. The 8-Cell 3D Macro-Cube
    wire [7:0] cell_out;
    
    geometric_qubit_virtualizer C0_Cell (.enable_phi_st(phase_field_active), .entanglement_in(1'b1), .q_state_out(cell_out[0]));
    assign master_entanglement_bus = cell_out[0]; 

    geometric_qubit_virtualizer C1_Cell (.enable_phi_st(phase_field_active), .entanglement_in(master_entanglement_bus), .q_state_out(cell_out[1]));
    geometric_qubit_virtualizer C2_Cell (.enable_phi_st(phase_field_active), .entanglement_in(master_entanglement_bus), .q_state_out(cell_out[2]));
    geometric_qubit_virtualizer C3_Cell (.enable_phi_st(phase_field_active), .entanglement_in(master_entanglement_bus), .q_state_out(cell_out[3]));
    geometric_qubit_virtualizer C4_Cell (.enable_phi_st(phase_field_active), .entanglement_in(master_entanglement_bus), .q_state_out(cell_out[4]));
    geometric_qubit_virtualizer C5_Cell (.enable_phi_st(phase_field_active), .entanglement_in(master_entanglement_bus), .q_state_out(cell_out[5]));
    geometric_qubit_virtualizer C6_Cell (.enable_phi_st(phase_field_active), .entanglement_in(master_entanglement_bus), .q_state_out(cell_out[6]));
    geometric_qubit_virtualizer C7_Cell (.enable_phi_st(phase_field_active), .entanglement_in(master_entanglement_bus), .q_state_out(cell_out[7]));

    // Map to LEDs
    assign led_0 = cell_out[0];
    assign led_1 = cell_out[1];
    assign led_2 = cell_out[2];
    assign led_3 = cell_out[3];

    // SPHY Core Integration
    wire [23:0] s_core_in = {12'b0, xadc_temp_data};
    wire [23:0] s_core_out;

    sphy_core engine (
        .clk(sys_clk),
        .rst_n(1'b1),
        .in_flux(s_core_in),
        .out(s_core_out)
    );

endmodule
