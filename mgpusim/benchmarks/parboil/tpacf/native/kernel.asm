	.text
	.hsa_code_object_version 2,1
	.hsa_code_object_isa 8,0,3,"AMD","AMDGPU"
	.protected	gen_hists       ; -- Begin function gen_hists
	.globl	gen_hists
	.p2align	8
	.type	gen_hists,@function
	.amdgpu_hsa_kernel gen_hists
gen_hists:                              ; @gen_hists
gen_hists$local:
	.amd_kernel_code_t
		amd_code_version_major = 1
		amd_code_version_minor = 2
		amd_machine_kind = 1
		amd_machine_version_major = 8
		amd_machine_version_minor = 0
		amd_machine_version_stepping = 3
		kernel_code_entry_byte_offset = 256
		kernel_code_prefetch_byte_size = 0
		granulated_workitem_vgpr_count = 3
		granulated_wavefront_sgpr_count = 5
		priority = 0
		float_mode = 192
		priv = 0
		enable_dx10_clamp = 1
		debug_mode = 0
		enable_ieee_mode = 1
		enable_wgp_mode = 0
		enable_mem_ordered = 0
		enable_fwd_progress = 0
		enable_sgpr_private_segment_wave_byte_offset = 0
		user_sgpr_count = 6
		enable_trap_handler = 0
		enable_sgpr_workgroup_id_x = 1
		enable_sgpr_workgroup_id_y = 0
		enable_sgpr_workgroup_id_z = 0
		enable_sgpr_workgroup_info = 0
		enable_vgpr_workitem_id = 0
		enable_exception_msb = 0
		granulated_lds_size = 0
		enable_exception = 0
		enable_sgpr_private_segment_buffer = 1
		enable_sgpr_dispatch_ptr = 0
		enable_sgpr_queue_ptr = 0
		enable_sgpr_kernarg_segment_ptr = 1
		enable_sgpr_dispatch_id = 0
		enable_sgpr_flat_scratch_init = 0
		enable_sgpr_private_segment_size = 0
		enable_sgpr_grid_workgroup_count_x = 0
		enable_sgpr_grid_workgroup_count_y = 0
		enable_sgpr_grid_workgroup_count_z = 0
		enable_wavefront_size32 = 0
		enable_ordered_append_gds = 0
		private_element_size = 1
		is_ptr64 = 1
		is_dynamic_callstack = 0
		is_debug_enabled = 0
		is_xnack_enabled = 0
		workitem_private_segment_byte_size = 0
		workgroup_group_segment_byte_size = 10240
		gds_segment_byte_size = 0
		kernarg_segment_byte_size = 88
		workgroup_fbarrier_count = 0
		wavefront_sgpr_count = 42
		workitem_vgpr_count = 15
		reserved_vgpr_first = 0
		reserved_vgpr_count = 0
		reserved_sgpr_first = 0
		reserved_sgpr_count = 0
		debug_wavefront_private_segment_offset_sgpr = 0
		debug_private_segment_buffer_sgpr = 0
		kernarg_segment_alignment = 4
		group_segment_alignment = 4
		private_segment_alignment = 4
		wavefront_size = 6
		call_convention = -1
		runtime_loader_kernel_symbol = 0
	.end_amd_kernel_code_t
; %bb.0:
	v_lshlrev_b32_e32 v3, 2, v0
	v_or_b32_e32 v6, 0x400, v3
	v_and_b32_e32 v2, 0x7f, v0
	v_and_b32_e32 v4, 0x200, v3
	v_lshlrev_b32_e32 v5, 2, v2
	v_and_b32_e32 v6, 0x600, v6
	v_mov_b32_e32 v1, 0
	v_or_b32_e32 v4, v4, v5
	s_mov_b32 m0, -1
	v_or_b32_e32 v6, v6, v5
	s_load_dwordx4 s[8:11], s[4:5], 0x0
	s_load_dwordx4 s[12:15], s[4:5], 0x10
	s_load_dwordx2 s[2:3], s[4:5], 0x18
	ds_write_b32 v4, v1
	ds_write_b32 v6, v1
	v_or_b32_e32 v6, 0x800, v3
	s_movk_i32 s0, 0xa00
	v_and_b32_e32 v6, s0, v6
	v_or_b32_e32 v6, v6, v5
	ds_write_b32 v6, v1
	v_or_b32_e32 v6, 0xc00, v3
	v_and_b32_e32 v6, 0xe00, v6
	v_or_b32_e32 v6, v6, v5
	ds_write_b32 v6, v1
	ds_write_b32 v4, v1 offset:4096
	v_or_b32_e32 v4, 0x1400, v3
	v_and_b32_e32 v4, 0x1600, v4
	v_or_b32_e32 v4, v4, v5
	ds_write_b32 v4, v1
	v_or_b32_e32 v4, 0x1800, v3
	v_or_b32_e32 v3, 0x1c00, v3
	v_and_b32_e32 v4, 0x1a00, v4
	v_and_b32_e32 v3, 0x1e00, v3
	v_or_b32_e32 v4, v4, v5
	v_or_b32_e32 v3, v3, v5
	ds_write_b32 v4, v1
	ds_write_b32 v3, v1
	v_or_b32_e32 v3, 0x800, v0
	v_cmp_gt_u32_e32 vcc, s0, v3
	s_and_saveexec_b64 s[0:1], vcc
; %bb.1:
	v_lshlrev_b32_e32 v3, 2, v3
	v_and_b32_e32 v3, 0x3e00, v3
	v_lshlrev_b32_e32 v4, 2, v2
	v_or_b32_e32 v3, v3, v4
	v_mov_b32_e32 v4, 0
	ds_write_b32 v3, v4
; %bb.2:
	s_or_b64 exec, exec, s[0:1]
	s_waitcnt lgkmcnt(0)
	s_add_i32 s15, s2, 1
	s_mul_i32 s0, s15, s3
	s_ashr_i32 s1, s0, 31
	s_lshl_b64 s[0:1], s[0:1], 2
	s_add_u32 s4, s10, s0
	s_addc_u32 s5, s11, s1
	v_lshlrev_b32_e32 v3, 2, v0
	v_or_b32_e32 v3, 0x2400, v3
	s_add_u32 s7, s4, s0
	s_addc_u32 s14, s5, s1
	v_and_b32_e32 v3, 0x2600, v3
	v_lshlrev_b32_e32 v2, 2, v2
	v_or_b32_e32 v2, v3, v2
	v_mov_b32_e32 v3, 0
	s_cmp_le_u32 s15, s6
	s_mov_b64 s[0:1], -1
	ds_write_b32 v2, v3
	s_cbranch_scc0 BB0_40
; %bb.3:                                ; %.thread
	s_cmp_eq_u32 s3, 0
	s_mov_b32 s17, 0
	s_cbranch_scc1 BB0_39
; %bb.4:                                ; %.split18.preheader
	s_sub_i32 s0, s6, s2
	s_mul_i32 s16, s0, s3
	s_lshl_b64 s[0:1], s[16:17], 2
	s_add_u32 s2, s7, s0
	s_addc_u32 s15, s14, s1
	s_add_u32 s18, s4, s0
	s_addc_u32 s19, s5, s1
	s_add_u32 s20, s10, s0
	s_addc_u32 s21, s11, s1
	s_add_i32 s0, s3, -1
	s_and_b32 s22, s3, 3
	v_lshrrev_b32_e32 v2, 1, v0
	v_cmp_gt_u32_e64 s[24:25], s0, 2
	s_and_b32 s23, s3, -4
	v_cmp_ne_u32_e64 s[26:27], s22, 0
	s_movk_i32 s28, 0xfe00
	s_mov_b32 s29, 0
                                        ; implicit-def: $vgpr5
                                        ; implicit-def: $vgpr4
                                        ; implicit-def: $vgpr3
	s_branch BB0_6
BB0_5:                                  ; %.us-lcssa
                                        ;   in Loop: Header=BB0_6 Depth=1
	s_addk_i32 s29, 0x100
	s_cmp_ge_u32 s29, s3
	s_cbranch_scc1 BB0_39
BB0_6:                                  ; %.split18
                                        ; =>This Loop Header: Depth=1
                                        ;     Child Loop BB0_11 Depth 2
                                        ;       Child Loop BB0_12 Depth 3
                                        ;       Child Loop BB0_17 Depth 3
                                        ;       Child Loop BB0_22 Depth 3
                                        ;       Child Loop BB0_27 Depth 3
                                        ;     Child Loop BB0_34 Depth 2
                                        ;       Child Loop BB0_35 Depth 3
	v_add_u32_e32 v6, vcc, s29, v0
	v_cmp_gt_u32_e32 vcc, s3, v6
	s_and_saveexec_b64 s[30:31], vcc
	s_cbranch_execz BB0_8
; %bb.7:                                ;   in Loop: Header=BB0_6 Depth=1
	v_mov_b32_e32 v7, 0
	s_waitcnt vmcnt(1) lgkmcnt(1)
	v_lshlrev_b64 v[3:4], 2, v[6:7]
	v_mov_b32_e32 v6, s21
	s_waitcnt vmcnt(0) lgkmcnt(0)
	v_add_u32_e64 v5, s[0:1], s20, v3
	v_addc_u32_e64 v6, s[0:1], v6, v4, s[0:1]
	v_mov_b32_e32 v8, s19
	v_add_u32_e64 v7, s[0:1], s18, v3
	v_addc_u32_e64 v8, s[0:1], v8, v4, s[0:1]
	v_mov_b32_e32 v10, s15
	v_add_u32_e64 v9, s[0:1], s2, v3
	v_addc_u32_e64 v10, s[0:1], v10, v4, s[0:1]
	flat_load_dword v3, v[5:6]
	flat_load_dword v4, v[7:8]
	flat_load_dword v5, v[9:10]
BB0_8:                                  ; %..split_crit_edge
                                        ;   in Loop: Header=BB0_6 Depth=1
	s_or_b64 exec, exec, s[30:31]
	s_xor_b64 s[0:1], vcc, -1
	s_andn2_b64 vcc, exec, s[24:25]
	s_mov_b32 s30, 0
	s_cbranch_vccnz BB0_32
; %bb.9:                                ; %..split_crit_edge.new.preheader
                                        ;   in Loop: Header=BB0_6 Depth=1
	s_mov_b32 s16, 0
	s_mov_b32 s30, s23
	s_branch BB0_11
BB0_10:                                 ;   in Loop: Header=BB0_11 Depth=2
	s_or_b64 exec, exec, s[34:35]
	s_add_i32 s16, s16, 4
	s_add_i32 s30, s30, -4
	s_cmp_eq_u32 s30, 0
	s_cbranch_scc1 BB0_31
BB0_11:                                 ; %..split_crit_edge.new
                                        ;   Parent Loop BB0_6 Depth=1
                                        ; =>  This Loop Header: Depth=2
                                        ;       Child Loop BB0_12 Depth 3
                                        ;       Child Loop BB0_17 Depth 3
                                        ;       Child Loop BB0_22 Depth 3
                                        ;       Child Loop BB0_27 Depth 3
	s_lshl_b64 s[34:35], s[16:17], 2
	s_add_u32 s36, s10, s34
	s_addc_u32 s37, s11, s35
	s_add_u32 s38, s4, s34
	s_addc_u32 s39, s5, s35
	s_add_u32 s34, s7, s34
	s_addc_u32 s35, s14, s35
	s_load_dword s31, s[36:37], 0x0
	s_load_dword s33, s[38:39], 0x0
	s_load_dword s34, s[34:35], 0x0
	v_mov_b32_e32 v6, 20
	v_mov_b32_e32 v9, 0
	s_waitcnt vmcnt(1) lgkmcnt(0)
	v_mul_f32_e32 v8, s33, v4
	v_mac_f32_e32 v8, s31, v3
	s_waitcnt vmcnt(0)
	v_mac_f32_e32 v8, s34, v5
	s_mov_b64 s[34:35], 0
BB0_12:                                 ;   Parent Loop BB0_6 Depth=1
                                        ;     Parent Loop BB0_11 Depth=2
                                        ; =>    This Inner Loop Header: Depth=3
	v_add_u32_e32 v7, vcc, v9, v6
	v_mov_b32_e32 v11, 0
	v_lshrrev_b32_e32 v10, 1, v7
	v_lshlrev_b64 v[11:12], 2, v[10:11]
	v_mov_b32_e32 v13, s13
	v_add_u32_e32 v11, vcc, s12, v11
	v_addc_u32_e32 v12, vcc, v13, v12, vcc
	flat_load_dword v7, v[11:12]
	s_waitcnt vmcnt(0) lgkmcnt(0)
	v_cmp_nge_f32_e32 vcc, v8, v7
	v_cndmask_b32_e32 v9, v9, v10, vcc
	v_cndmask_b32_e32 v6, v10, v6, vcc
	v_add_u32_e32 v7, vcc, 1, v9
	v_cmp_le_u32_e32 vcc, v6, v7
	s_or_b64 s[34:35], vcc, s[34:35]
	s_andn2_b64 exec, exec, s[34:35]
	s_cbranch_execnz BB0_12
; %bb.13:                               ;   in Loop: Header=BB0_11 Depth=2
	s_or_b64 exec, exec, s[34:35]
	v_mov_b32_e32 v10, 0
	v_lshlrev_b64 v[9:10], 2, v[9:10]
	v_mov_b32_e32 v7, s13
	v_add_u32_e32 v9, vcc, s12, v9
	v_addc_u32_e32 v10, vcc, v7, v10, vcc
	flat_load_dword v7, v[9:10]
	s_waitcnt vmcnt(0) lgkmcnt(0)
	v_cmp_lt_f32_e32 vcc, v8, v7
	s_and_saveexec_b64 s[34:35], vcc
	s_cbranch_execz BB0_16
; %bb.14:                               ;   in Loop: Header=BB0_11 Depth=2
	v_mov_b32_e32 v7, 0
	v_lshlrev_b64 v[9:10], 2, v[6:7]
	v_mov_b32_e32 v7, s13
	v_add_u32_e32 v9, vcc, s12, v9
	v_addc_u32_e32 v10, vcc, v7, v10, vcc
	flat_load_dword v7, v[9:10]
	s_xor_b64 s[36:37], s[0:1], -1
	s_waitcnt vmcnt(0) lgkmcnt(0)
	v_cmp_ge_f32_e32 vcc, v8, v7
	s_and_b64 s[36:37], vcc, s[36:37]
	s_and_b64 exec, exec, s[36:37]
; %bb.15:                               ;   in Loop: Header=BB0_11 Depth=2
	v_lshlrev_b32_e32 v6, 9, v6
	v_lshlrev_b32_e32 v7, 2, v2
	v_or_b32_e32 v6, v6, v7
	v_add_u32_e32 v6, vcc, s28, v6
	v_mov_b32_e32 v7, 1
	ds_add_u32 v6, v7
BB0_16:                                 ;   in Loop: Header=BB0_11 Depth=2
	s_or_b64 exec, exec, s[34:35]
	s_or_b32 s34, s16, 1
	s_mov_b32 s35, s17
	s_lshl_b64 s[34:35], s[34:35], 2
	s_add_u32 s36, s10, s34
	s_addc_u32 s37, s11, s35
	s_add_u32 s38, s4, s34
	s_addc_u32 s39, s5, s35
	s_add_u32 s34, s7, s34
	s_addc_u32 s35, s14, s35
	s_load_dword s31, s[36:37], 0x0
	s_load_dword s33, s[38:39], 0x0
	s_load_dword s34, s[34:35], 0x0
	v_mov_b32_e32 v6, 20
	v_mov_b32_e32 v9, 0
	s_waitcnt lgkmcnt(0)
	v_mul_f32_e32 v8, s33, v4
	v_mac_f32_e32 v8, s31, v3
	v_mac_f32_e32 v8, s34, v5
	s_mov_b64 s[34:35], 0
BB0_17:                                 ;   Parent Loop BB0_6 Depth=1
                                        ;     Parent Loop BB0_11 Depth=2
                                        ; =>    This Inner Loop Header: Depth=3
	v_add_u32_e32 v7, vcc, v9, v6
	v_mov_b32_e32 v11, 0
	v_lshrrev_b32_e32 v10, 1, v7
	v_lshlrev_b64 v[11:12], 2, v[10:11]
	v_mov_b32_e32 v13, s13
	v_add_u32_e32 v11, vcc, s12, v11
	v_addc_u32_e32 v12, vcc, v13, v12, vcc
	flat_load_dword v7, v[11:12]
	s_waitcnt vmcnt(0) lgkmcnt(0)
	v_cmp_nge_f32_e32 vcc, v8, v7
	v_cndmask_b32_e32 v9, v9, v10, vcc
	v_cndmask_b32_e32 v6, v10, v6, vcc
	v_add_u32_e32 v7, vcc, 1, v9
	v_cmp_le_u32_e32 vcc, v6, v7
	s_or_b64 s[34:35], vcc, s[34:35]
	s_andn2_b64 exec, exec, s[34:35]
	s_cbranch_execnz BB0_17
; %bb.18:                               ;   in Loop: Header=BB0_11 Depth=2
	s_or_b64 exec, exec, s[34:35]
	v_mov_b32_e32 v10, 0
	v_lshlrev_b64 v[9:10], 2, v[9:10]
	v_mov_b32_e32 v7, s13
	v_add_u32_e32 v9, vcc, s12, v9
	v_addc_u32_e32 v10, vcc, v7, v10, vcc
	flat_load_dword v7, v[9:10]
	s_waitcnt vmcnt(0) lgkmcnt(0)
	v_cmp_lt_f32_e32 vcc, v8, v7
	s_and_saveexec_b64 s[34:35], vcc
	s_cbranch_execz BB0_21
; %bb.19:                               ;   in Loop: Header=BB0_11 Depth=2
	v_mov_b32_e32 v7, 0
	v_lshlrev_b64 v[9:10], 2, v[6:7]
	v_mov_b32_e32 v7, s13
	v_add_u32_e32 v9, vcc, s12, v9
	v_addc_u32_e32 v10, vcc, v7, v10, vcc
	flat_load_dword v7, v[9:10]
	s_xor_b64 s[36:37], s[0:1], -1
	s_waitcnt vmcnt(0) lgkmcnt(0)
	v_cmp_ge_f32_e32 vcc, v8, v7
	s_and_b64 s[36:37], vcc, s[36:37]
	s_and_b64 exec, exec, s[36:37]
; %bb.20:                               ;   in Loop: Header=BB0_11 Depth=2
	v_lshlrev_b32_e32 v6, 9, v6
	v_lshlrev_b32_e32 v7, 2, v2
	v_or_b32_e32 v6, v6, v7
	v_add_u32_e32 v6, vcc, s28, v6
	v_mov_b32_e32 v7, 1
	ds_add_u32 v6, v7
BB0_21:                                 ;   in Loop: Header=BB0_11 Depth=2
	s_or_b64 exec, exec, s[34:35]
	s_or_b32 s34, s16, 2
	s_mov_b32 s35, s17
	s_lshl_b64 s[34:35], s[34:35], 2
	s_add_u32 s36, s10, s34
	s_addc_u32 s37, s11, s35
	s_add_u32 s38, s4, s34
	s_addc_u32 s39, s5, s35
	s_add_u32 s34, s7, s34
	s_addc_u32 s35, s14, s35
	s_load_dword s31, s[36:37], 0x0
	s_load_dword s33, s[38:39], 0x0
	s_load_dword s34, s[34:35], 0x0
	v_mov_b32_e32 v6, 20
	v_mov_b32_e32 v9, 0
	s_waitcnt lgkmcnt(0)
	v_mul_f32_e32 v8, s33, v4
	v_mac_f32_e32 v8, s31, v3
	v_mac_f32_e32 v8, s34, v5
	s_mov_b64 s[34:35], 0
BB0_22:                                 ;   Parent Loop BB0_6 Depth=1
                                        ;     Parent Loop BB0_11 Depth=2
                                        ; =>    This Inner Loop Header: Depth=3
	v_add_u32_e32 v7, vcc, v9, v6
	v_mov_b32_e32 v11, 0
	v_lshrrev_b32_e32 v10, 1, v7
	v_lshlrev_b64 v[11:12], 2, v[10:11]
	v_mov_b32_e32 v13, s13
	v_add_u32_e32 v11, vcc, s12, v11
	v_addc_u32_e32 v12, vcc, v13, v12, vcc
	flat_load_dword v7, v[11:12]
	s_waitcnt vmcnt(0) lgkmcnt(0)
	v_cmp_nge_f32_e32 vcc, v8, v7
	v_cndmask_b32_e32 v9, v9, v10, vcc
	v_cndmask_b32_e32 v6, v10, v6, vcc
	v_add_u32_e32 v7, vcc, 1, v9
	v_cmp_le_u32_e32 vcc, v6, v7
	s_or_b64 s[34:35], vcc, s[34:35]
	s_andn2_b64 exec, exec, s[34:35]
	s_cbranch_execnz BB0_22
; %bb.23:                               ;   in Loop: Header=BB0_11 Depth=2
	s_or_b64 exec, exec, s[34:35]
	v_mov_b32_e32 v10, 0
	v_lshlrev_b64 v[9:10], 2, v[9:10]
	v_mov_b32_e32 v7, s13
	v_add_u32_e32 v9, vcc, s12, v9
	v_addc_u32_e32 v10, vcc, v7, v10, vcc
	flat_load_dword v7, v[9:10]
	s_waitcnt vmcnt(0) lgkmcnt(0)
	v_cmp_lt_f32_e32 vcc, v8, v7
	s_and_saveexec_b64 s[34:35], vcc
	s_cbranch_execz BB0_26
; %bb.24:                               ;   in Loop: Header=BB0_11 Depth=2
	v_mov_b32_e32 v7, 0
	v_lshlrev_b64 v[9:10], 2, v[6:7]
	v_mov_b32_e32 v7, s13
	v_add_u32_e32 v9, vcc, s12, v9
	v_addc_u32_e32 v10, vcc, v7, v10, vcc
	flat_load_dword v7, v[9:10]
	s_xor_b64 s[36:37], s[0:1], -1
	s_waitcnt vmcnt(0) lgkmcnt(0)
	v_cmp_ge_f32_e32 vcc, v8, v7
	s_and_b64 s[36:37], vcc, s[36:37]
	s_and_b64 exec, exec, s[36:37]
; %bb.25:                               ;   in Loop: Header=BB0_11 Depth=2
	v_lshlrev_b32_e32 v6, 9, v6
	v_lshlrev_b32_e32 v7, 2, v2
	v_or_b32_e32 v6, v6, v7
	v_add_u32_e32 v6, vcc, s28, v6
	v_mov_b32_e32 v7, 1
	ds_add_u32 v6, v7
BB0_26:                                 ;   in Loop: Header=BB0_11 Depth=2
	s_or_b64 exec, exec, s[34:35]
	s_or_b32 s34, s16, 3
	s_mov_b32 s35, s17
	s_lshl_b64 s[34:35], s[34:35], 2
	s_add_u32 s36, s10, s34
	s_addc_u32 s37, s11, s35
	s_add_u32 s38, s4, s34
	s_addc_u32 s39, s5, s35
	s_add_u32 s34, s7, s34
	s_addc_u32 s35, s14, s35
	s_load_dword s31, s[36:37], 0x0
	s_load_dword s33, s[38:39], 0x0
	s_load_dword s34, s[34:35], 0x0
	v_mov_b32_e32 v6, 20
	v_mov_b32_e32 v9, 0
	s_waitcnt lgkmcnt(0)
	v_mul_f32_e32 v8, s33, v4
	v_mac_f32_e32 v8, s31, v3
	v_mac_f32_e32 v8, s34, v5
	s_mov_b64 s[34:35], 0
BB0_27:                                 ;   Parent Loop BB0_6 Depth=1
                                        ;     Parent Loop BB0_11 Depth=2
                                        ; =>    This Inner Loop Header: Depth=3
	v_add_u32_e32 v7, vcc, v9, v6
	v_mov_b32_e32 v11, 0
	v_lshrrev_b32_e32 v10, 1, v7
	v_lshlrev_b64 v[11:12], 2, v[10:11]
	v_mov_b32_e32 v13, s13
	v_add_u32_e32 v11, vcc, s12, v11
	v_addc_u32_e32 v12, vcc, v13, v12, vcc
	flat_load_dword v7, v[11:12]
	s_waitcnt vmcnt(0) lgkmcnt(0)
	v_cmp_nge_f32_e32 vcc, v8, v7
	v_cndmask_b32_e32 v9, v9, v10, vcc
	v_cndmask_b32_e32 v6, v10, v6, vcc
	v_add_u32_e32 v7, vcc, 1, v9
	v_cmp_le_u32_e32 vcc, v6, v7
	s_or_b64 s[34:35], vcc, s[34:35]
	s_andn2_b64 exec, exec, s[34:35]
	s_cbranch_execnz BB0_27
; %bb.28:                               ;   in Loop: Header=BB0_11 Depth=2
	s_or_b64 exec, exec, s[34:35]
	v_mov_b32_e32 v10, 0
	v_lshlrev_b64 v[9:10], 2, v[9:10]
	v_mov_b32_e32 v7, s13
	v_add_u32_e32 v9, vcc, s12, v9
	v_addc_u32_e32 v10, vcc, v7, v10, vcc
	flat_load_dword v7, v[9:10]
	s_waitcnt vmcnt(0) lgkmcnt(0)
	v_cmp_lt_f32_e32 vcc, v8, v7
	s_and_saveexec_b64 s[34:35], vcc
	s_cbranch_execz BB0_10
; %bb.29:                               ;   in Loop: Header=BB0_11 Depth=2
	v_mov_b32_e32 v7, 0
	v_lshlrev_b64 v[9:10], 2, v[6:7]
	v_mov_b32_e32 v7, s13
	v_add_u32_e32 v9, vcc, s12, v9
	v_addc_u32_e32 v10, vcc, v7, v10, vcc
	flat_load_dword v7, v[9:10]
	s_xor_b64 s[36:37], s[0:1], -1
	s_waitcnt vmcnt(0) lgkmcnt(0)
	v_cmp_ge_f32_e32 vcc, v8, v7
	s_and_b64 s[36:37], vcc, s[36:37]
	s_and_b64 exec, exec, s[36:37]
	s_cbranch_execz BB0_10
; %bb.30:                               ;   in Loop: Header=BB0_11 Depth=2
	v_lshlrev_b32_e32 v6, 9, v6
	v_lshlrev_b32_e32 v7, 2, v2
	v_or_b32_e32 v6, v6, v7
	v_add_u32_e32 v6, vcc, s28, v6
	v_mov_b32_e32 v7, 1
	ds_add_u32 v6, v7
	s_branch BB0_10
BB0_31:                                 ; %Flow152
                                        ;   in Loop: Header=BB0_6 Depth=1
	s_mov_b32 s30, s16
BB0_32:                                 ; %.us-lcssa.unr-lcssa
                                        ;   in Loop: Header=BB0_6 Depth=1
	s_andn2_b64 vcc, exec, s[26:27]
	s_mov_b32 s16, s22
	s_cbranch_vccz BB0_34
	s_branch BB0_5
BB0_33:                                 ;   in Loop: Header=BB0_34 Depth=2
	s_or_b64 exec, exec, s[34:35]
	s_add_i32 s30, s30, 1
	s_add_i32 s16, s16, -1
	s_cmp_lg_u32 s16, 0
	s_cbranch_scc0 BB0_5
BB0_34:                                 ; %.epil.preheader
                                        ;   Parent Loop BB0_6 Depth=1
                                        ; =>  This Loop Header: Depth=2
                                        ;       Child Loop BB0_35 Depth 3
	s_mov_b32 s31, s17
	s_lshl_b64 s[34:35], s[30:31], 2
	s_add_u32 s36, s10, s34
	s_addc_u32 s37, s11, s35
	s_add_u32 s38, s4, s34
	s_addc_u32 s39, s5, s35
	s_add_u32 s34, s7, s34
	s_addc_u32 s35, s14, s35
	s_load_dword s31, s[36:37], 0x0
	s_load_dword s33, s[38:39], 0x0
	s_load_dword s34, s[34:35], 0x0
	v_mov_b32_e32 v6, 20
	v_mov_b32_e32 v9, 0
	s_waitcnt vmcnt(1) lgkmcnt(0)
	v_mul_f32_e32 v8, s33, v4
	v_mac_f32_e32 v8, s31, v3
	s_waitcnt vmcnt(0)
	v_mac_f32_e32 v8, s34, v5
	s_mov_b64 s[34:35], 0
BB0_35:                                 ;   Parent Loop BB0_6 Depth=1
                                        ;     Parent Loop BB0_34 Depth=2
                                        ; =>    This Inner Loop Header: Depth=3
	v_add_u32_e32 v7, vcc, v9, v6
	v_mov_b32_e32 v11, 0
	v_lshrrev_b32_e32 v10, 1, v7
	v_lshlrev_b64 v[11:12], 2, v[10:11]
	v_mov_b32_e32 v13, s13
	v_add_u32_e32 v11, vcc, s12, v11
	v_addc_u32_e32 v12, vcc, v13, v12, vcc
	flat_load_dword v7, v[11:12]
	s_waitcnt vmcnt(0) lgkmcnt(0)
	v_cmp_nge_f32_e32 vcc, v8, v7
	v_cndmask_b32_e32 v9, v9, v10, vcc
	v_cndmask_b32_e32 v6, v10, v6, vcc
	v_add_u32_e32 v7, vcc, 1, v9
	v_cmp_le_u32_e32 vcc, v6, v7
	s_or_b64 s[34:35], vcc, s[34:35]
	s_andn2_b64 exec, exec, s[34:35]
	s_cbranch_execnz BB0_35
; %bb.36:                               ;   in Loop: Header=BB0_34 Depth=2
	s_or_b64 exec, exec, s[34:35]
	v_mov_b32_e32 v10, 0
	v_lshlrev_b64 v[9:10], 2, v[9:10]
	v_mov_b32_e32 v7, s13
	v_add_u32_e32 v9, vcc, s12, v9
	v_addc_u32_e32 v10, vcc, v7, v10, vcc
	flat_load_dword v7, v[9:10]
	s_waitcnt vmcnt(0) lgkmcnt(0)
	v_cmp_lt_f32_e32 vcc, v8, v7
	s_and_saveexec_b64 s[34:35], vcc
	s_cbranch_execz BB0_33
; %bb.37:                               ;   in Loop: Header=BB0_34 Depth=2
	v_mov_b32_e32 v7, 0
	v_lshlrev_b64 v[9:10], 2, v[6:7]
	v_mov_b32_e32 v7, s13
	v_add_u32_e32 v9, vcc, s12, v9
	v_addc_u32_e32 v10, vcc, v7, v10, vcc
	flat_load_dword v7, v[9:10]
	s_xor_b64 s[36:37], s[0:1], -1
	s_waitcnt vmcnt(0) lgkmcnt(0)
	v_cmp_ge_f32_e32 vcc, v8, v7
	s_and_b64 s[36:37], vcc, s[36:37]
	s_and_b64 exec, exec, s[36:37]
	s_cbranch_execz BB0_33
; %bb.38:                               ;   in Loop: Header=BB0_34 Depth=2
	v_lshlrev_b32_e32 v6, 9, v6
	v_lshlrev_b32_e32 v7, 2, v2
	v_or_b32_e32 v6, v6, v7
	v_add_u32_e32 v6, vcc, s28, v6
	v_mov_b32_e32 v7, 1
	ds_add_u32 v6, v7
	s_branch BB0_33
BB0_39:                                 ; %Flow155
	s_mov_b64 s[0:1], 0
BB0_40:                                 ; %Flow165
	s_andn2_b64 vcc, exec, s[0:1]
	s_cbranch_vccnz BB0_61
; %bb.41:
	s_cmp_lg_u32 s3, 0
	s_mov_b32 s17, 0
	s_cbranch_scc0 BB0_61
; %bb.42:                               ; %.split18.us.preheader
	s_mul_i32 s16, s6, s3
	s_lshl_b64 s[0:1], s[16:17], 2
	s_add_u32 s2, s10, s0
	s_addc_u32 s10, s11, s1
	s_add_u32 s4, s4, s0
	s_addc_u32 s5, s5, s1
	s_add_u32 s7, s7, s0
	s_addc_u32 s11, s14, s1
	v_lshrrev_b32_e32 v2, 1, v0
	s_add_i32 s14, s3, -1
	s_movk_i32 s15, 0x100
	s_movk_i32 s18, 0x100
	s_mov_b32 s19, 0
	s_waitcnt vmcnt(2) lgkmcnt(2)
                                        ; implicit-def: $vgpr3
	s_waitcnt vmcnt(1) lgkmcnt(1)
                                        ; implicit-def: $vgpr4
	s_waitcnt vmcnt(0) lgkmcnt(0)
                                        ; implicit-def: $vgpr5
	s_branch BB0_44
BB0_43:                                 ; %Flow162
                                        ;   in Loop: Header=BB0_44 Depth=1
	s_or_b64 exec, exec, s[20:21]
	s_add_i32 s19, s19, s15
	s_add_i32 s18, s18, s15
	s_cmp_ge_u32 s19, s3
	s_cbranch_scc1 BB0_61
BB0_44:                                 ; %.split18.us
                                        ; =>This Loop Header: Depth=1
                                        ;     Child Loop BB0_47 Depth 2
                                        ;       Child Loop BB0_49 Depth 3
                                        ;     Child Loop BB0_55 Depth 2
                                        ;       Child Loop BB0_57 Depth 3
	v_add_u32_e32 v6, vcc, s19, v0
	v_cmp_le_u32_e32 vcc, s3, v6
	s_and_saveexec_b64 s[0:1], vcc
	s_xor_b64 s[20:21], exec, s[0:1]
	s_cbranch_execz BB0_51
; %bb.45:                               ; %.split.us.us.split.us.preheader
                                        ;   in Loop: Header=BB0_44 Depth=1
	v_mov_b32_e32 v7, s14
	v_cmp_gt_u32_e64 s[0:1], s18, v7
	v_mov_b32_e32 v7, s18
	v_cmp_le_u32_e64 s[22:23], s14, v7
	s_mov_b32 s24, 0
	s_branch BB0_47
BB0_46:                                 ; %Flow156
                                        ;   in Loop: Header=BB0_47 Depth=2
	s_and_b64 vcc, exec, s[26:27]
	s_cbranch_vccnz BB0_51
BB0_47:                                 ; %.split.us.us.split.us
                                        ;   Parent Loop BB0_44 Depth=1
                                        ; =>  This Loop Header: Depth=2
                                        ;       Child Loop BB0_49 Depth 3
	s_andn2_b64 vcc, exec, s[0:1]
	s_mov_b64 s[26:27], -1
	s_cbranch_vccnz BB0_46
; %bb.48:                               ;   in Loop: Header=BB0_47 Depth=2
	s_mov_b32 s25, s17
	s_lshl_b64 s[26:27], s[24:25], 2
	s_add_u32 s28, s2, s26
	s_addc_u32 s29, s10, s27
	s_add_u32 s30, s4, s26
	s_addc_u32 s31, s5, s27
	s_add_u32 s26, s7, s26
	s_addc_u32 s27, s11, s27
	s_load_dword s16, s[28:29], 0x0
	s_load_dword s25, s[30:31], 0x0
	s_load_dword s26, s[26:27], 0x0
	v_mov_b32_e32 v8, 20
	v_mov_b32_e32 v9, 0
	s_waitcnt vmcnt(1) lgkmcnt(0)
	v_mul_f32_e32 v7, s25, v4
	v_mac_f32_e32 v7, s16, v5
	s_waitcnt vmcnt(0)
	v_mac_f32_e32 v7, s26, v3
	s_mov_b64 s[26:27], 0
BB0_49:                                 ;   Parent Loop BB0_44 Depth=1
                                        ;     Parent Loop BB0_47 Depth=2
                                        ; =>    This Inner Loop Header: Depth=3
	v_add_u32_e32 v10, vcc, v9, v8
	v_mov_b32_e32 v11, 0
	v_lshrrev_b32_e32 v10, 1, v10
	v_lshlrev_b64 v[11:12], 2, v[10:11]
	v_mov_b32_e32 v13, s13
	v_add_u32_e32 v11, vcc, s12, v11
	v_addc_u32_e32 v12, vcc, v13, v12, vcc
	flat_load_dword v11, v[11:12]
	s_waitcnt vmcnt(0) lgkmcnt(0)
	v_cmp_nge_f32_e32 vcc, v7, v11
	v_cndmask_b32_e32 v9, v9, v10, vcc
	v_cndmask_b32_e32 v8, v10, v8, vcc
	v_add_u32_e32 v10, vcc, 1, v9
	v_cmp_le_u32_e32 vcc, v8, v10
	s_or_b64 s[26:27], vcc, s[26:27]
	s_andn2_b64 exec, exec, s[26:27]
	s_cbranch_execnz BB0_49
; %bb.50:                               ;   in Loop: Header=BB0_47 Depth=2
	s_or_b64 exec, exec, s[26:27]
	s_add_i32 s24, s24, 1
	s_mov_b64 s[26:27], s[22:23]
	s_branch BB0_46
BB0_51:                                 ; %Flow161
                                        ;   in Loop: Header=BB0_44 Depth=1
	s_or_saveexec_b64 s[20:21], s[20:21]
	s_xor_b64 exec, exec, s[20:21]
	s_cbranch_execz BB0_43
; %bb.52:                               ; %.split.us.us.split.preheader
                                        ;   in Loop: Header=BB0_44 Depth=1
	v_mov_b32_e32 v7, 0
	s_waitcnt vmcnt(0) lgkmcnt(0)
	v_lshlrev_b64 v[3:4], 2, v[6:7]
	v_mov_b32_e32 v5, s10
	v_add_u32_e32 v7, vcc, s2, v3
	v_addc_u32_e32 v8, vcc, v5, v4, vcc
	v_mov_b32_e32 v5, s5
	v_add_u32_e32 v9, vcc, s4, v3
	v_addc_u32_e32 v10, vcc, v5, v4, vcc
	v_mov_b32_e32 v5, s11
	v_add_u32_e32 v11, vcc, s7, v3
	v_addc_u32_e32 v12, vcc, v5, v4, vcc
	flat_load_dword v5, v[7:8]
	flat_load_dword v4, v[9:10]
	flat_load_dword v3, v[11:12]
	s_mov_b32 s16, 0
	s_branch BB0_55
BB0_53:                                 ;   in Loop: Header=BB0_55 Depth=2
	s_or_b64 exec, exec, s[22:23]
	s_add_i32 s16, s16, 1
	v_mov_b32_e32 v7, s3
	v_cmp_eq_u32_e64 s[0:1], s16, v7
BB0_54:                                 ; %Flow159
                                        ;   in Loop: Header=BB0_55 Depth=2
	s_and_b64 vcc, exec, s[0:1]
	s_cbranch_vccnz BB0_43
BB0_55:                                 ; %.split.us.us.split
                                        ;   Parent Loop BB0_44 Depth=1
                                        ; =>  This Loop Header: Depth=2
                                        ;       Child Loop BB0_57 Depth 3
	s_cmp_eq_u32 s16, s18
	s_mov_b64 s[0:1], -1
	s_cbranch_scc1 BB0_54
; %bb.56:                               ;   in Loop: Header=BB0_55 Depth=2
	s_lshl_b64 s[0:1], s[16:17], 2
	s_add_u32 s22, s2, s0
	s_addc_u32 s23, s10, s1
	s_add_u32 s24, s4, s0
	s_addc_u32 s25, s5, s1
	s_add_u32 s0, s7, s0
	s_addc_u32 s1, s11, s1
	s_load_dword s22, s[22:23], 0x0
	s_load_dword s23, s[24:25], 0x0
	s_load_dword s0, s[0:1], 0x0
	v_mov_b32_e32 v7, 20
	v_mov_b32_e32 v10, 0
	s_waitcnt vmcnt(1) lgkmcnt(0)
	v_mul_f32_e32 v9, s23, v4
	v_mac_f32_e32 v9, s22, v5
	s_waitcnt vmcnt(0)
	v_mac_f32_e32 v9, s0, v3
	s_mov_b64 s[0:1], 0
BB0_57:                                 ;   Parent Loop BB0_44 Depth=1
                                        ;     Parent Loop BB0_55 Depth=2
                                        ; =>    This Inner Loop Header: Depth=3
	v_add_u32_e32 v8, vcc, v10, v7
	v_mov_b32_e32 v12, 0
	v_lshrrev_b32_e32 v11, 1, v8
	v_lshlrev_b64 v[12:13], 2, v[11:12]
	v_mov_b32_e32 v14, s13
	v_add_u32_e32 v12, vcc, s12, v12
	v_addc_u32_e32 v13, vcc, v14, v13, vcc
	flat_load_dword v8, v[12:13]
	s_waitcnt vmcnt(0) lgkmcnt(0)
	v_cmp_nge_f32_e32 vcc, v9, v8
	v_cndmask_b32_e32 v10, v10, v11, vcc
	v_cndmask_b32_e32 v7, v11, v7, vcc
	v_add_u32_e32 v8, vcc, 1, v10
	v_cmp_le_u32_e32 vcc, v7, v8
	s_or_b64 s[0:1], vcc, s[0:1]
	s_andn2_b64 exec, exec, s[0:1]
	s_cbranch_execnz BB0_57
; %bb.58:                               ;   in Loop: Header=BB0_55 Depth=2
	s_or_b64 exec, exec, s[0:1]
	v_mov_b32_e32 v11, 0
	v_lshlrev_b64 v[10:11], 2, v[10:11]
	v_mov_b32_e32 v8, s13
	v_add_u32_e32 v10, vcc, s12, v10
	v_addc_u32_e32 v11, vcc, v8, v11, vcc
	flat_load_dword v8, v[10:11]
	s_waitcnt vmcnt(0) lgkmcnt(0)
	v_cmp_lt_f32_e32 vcc, v9, v8
	s_and_saveexec_b64 s[22:23], vcc
	s_cbranch_execz BB0_53
; %bb.59:                               ;   in Loop: Header=BB0_55 Depth=2
	v_mov_b32_e32 v8, 0
	v_lshlrev_b64 v[10:11], 2, v[7:8]
	v_mov_b32_e32 v8, s13
	v_add_u32_e32 v10, vcc, s12, v10
	v_addc_u32_e32 v11, vcc, v8, v11, vcc
	flat_load_dword v8, v[10:11]
	v_cmp_lt_u32_e64 s[0:1], s16, v6
	s_waitcnt vmcnt(0) lgkmcnt(0)
	v_cmp_ge_f32_e32 vcc, v9, v8
	s_and_b64 s[0:1], vcc, s[0:1]
	s_and_b64 exec, exec, s[0:1]
	s_cbranch_execz BB0_53
; %bb.60:                               ;   in Loop: Header=BB0_55 Depth=2
	v_lshlrev_b32_e32 v7, 9, v7
	v_lshlrev_b32_e32 v8, 2, v2
	v_or_b32_e32 v7, v7, v8
	v_add_u32_e32 v7, vcc, 0xfffffe00, v7
	v_mov_b32_e32 v8, 1
	ds_add_u32 v7, v8
	s_branch BB0_53
BB0_61:                                 ; %.loopexit
	s_waitcnt vmcnt(0) lgkmcnt(0)
	v_and_b32_e32 v3, 63, v0
	v_lshlrev_b32_e32 v2, 3, v0
	v_and_b32_e32 v2, 0x600, v2
	v_lshlrev_b32_e32 v4, 2, v3
	v_or_b32_e32 v2, v2, v4
	s_waitcnt lgkmcnt(0)
	s_barrier
	s_waitcnt lgkmcnt(0)
	ds_read2st64_b32 v[4:5], v2 offset1:1
	v_cmp_gt_u32_e64 s[0:1], 32, v3
	s_mov_b64 s[2:3], 0
	s_waitcnt lgkmcnt(0)
	v_add_u32_e32 v4, vcc, v5, v4
	ds_write_b32 v2, v4
	s_waitcnt lgkmcnt(0)
	s_barrier
	s_waitcnt lgkmcnt(0)
	ds_read2st64_b32 v[4:5], v2 offset0:8 offset1:9
	s_waitcnt lgkmcnt(0)
	v_add_u32_e32 v4, vcc, v5, v4
	ds_write_b32 v2, v4 offset:2048
	s_waitcnt lgkmcnt(0)
	s_barrier
	s_waitcnt lgkmcnt(0)
	ds_read2st64_b32 v[4:5], v2 offset0:16 offset1:17
	s_waitcnt lgkmcnt(0)
	v_add_u32_e32 v4, vcc, v5, v4
	ds_write_b32 v2, v4 offset:4096
	s_waitcnt lgkmcnt(0)
	s_barrier
	s_waitcnt lgkmcnt(0)
	ds_read2st64_b32 v[4:5], v2 offset0:24 offset1:25
	s_waitcnt lgkmcnt(0)
	v_add_u32_e32 v4, vcc, v5, v4
	ds_write_b32 v2, v4 offset:6144
	s_waitcnt lgkmcnt(0)
	s_barrier
	s_waitcnt lgkmcnt(0)
	ds_read2st64_b32 v[4:5], v2 offset0:32 offset1:33
	s_waitcnt lgkmcnt(0)
	v_add_u32_e32 v4, vcc, v5, v4
	ds_write_b32 v2, v4 offset:8192
	s_waitcnt lgkmcnt(0)
	s_barrier
	s_waitcnt lgkmcnt(0)
	s_and_saveexec_b64 s[4:5], s[0:1]
	s_cbranch_execz BB0_63
; %bb.62:
	ds_read2_b32 v[4:5], v2 offset1:32
	s_mov_b64 s[2:3], exec
	s_waitcnt lgkmcnt(0)
	v_add_u32_e32 v4, vcc, v5, v4
	ds_write_b32 v2, v4
BB0_63:
	s_or_b64 exec, exec, s[4:5]
	s_waitcnt lgkmcnt(0)
	s_barrier
	s_waitcnt lgkmcnt(0)
	s_and_saveexec_b64 s[4:5], s[2:3]
	s_cbranch_execz BB0_65
; %bb.64:
	s_movk_i32 s2, 0x800
	v_add_u32_e32 v4, vcc, s2, v2
	ds_read2_b32 v[4:5], v4 offset1:32
	s_waitcnt lgkmcnt(0)
	v_add_u32_e32 v4, vcc, v5, v4
	ds_write_b32 v2, v4 offset:2048
BB0_65:
	s_or_b64 exec, exec, s[4:5]
	s_mov_b64 s[2:3], 0
	s_waitcnt lgkmcnt(0)
	s_barrier
	s_waitcnt lgkmcnt(0)
	s_and_saveexec_b64 s[4:5], s[0:1]
	s_cbranch_execz BB0_67
; %bb.66:
	s_movk_i32 s7, 0x1000
	v_add_u32_e32 v4, vcc, s7, v2
	ds_read2_b32 v[4:5], v4 offset1:32
	s_mov_b64 s[2:3], exec
	s_waitcnt lgkmcnt(0)
	v_add_u32_e32 v4, vcc, v5, v4
	ds_write_b32 v2, v4 offset:4096
BB0_67:
	s_or_b64 exec, exec, s[4:5]
	s_waitcnt lgkmcnt(0)
	s_barrier
	s_waitcnt lgkmcnt(0)
	s_and_saveexec_b64 s[4:5], s[2:3]
	s_cbranch_execz BB0_69
; %bb.68:
	s_movk_i32 s2, 0x1800
	v_add_u32_e32 v4, vcc, s2, v2
	ds_read2_b32 v[4:5], v4 offset1:32
	s_waitcnt lgkmcnt(0)
	v_add_u32_e32 v4, vcc, v5, v4
	ds_write_b32 v2, v4 offset:6144
BB0_69:
	s_or_b64 exec, exec, s[4:5]
	s_waitcnt lgkmcnt(0)
	s_barrier
	s_waitcnt lgkmcnt(0)
	s_and_saveexec_b64 s[2:3], s[0:1]
	s_cbranch_execz BB0_71
; %bb.70:
	s_movk_i32 s0, 0x2000
	v_add_u32_e32 v4, vcc, s0, v2
	ds_read2_b32 v[4:5], v4 offset1:32
	s_waitcnt lgkmcnt(0)
	v_add_u32_e32 v4, vcc, v5, v4
	ds_write_b32 v2, v4 offset:8192
BB0_71:
	s_or_b64 exec, exec, s[2:3]
	v_cmp_gt_u32_e64 s[0:1], 16, v3
	s_mov_b64 s[2:3], 0
	s_waitcnt lgkmcnt(0)
	s_barrier
	s_waitcnt lgkmcnt(0)
	s_and_saveexec_b64 s[4:5], s[0:1]
	s_cbranch_execz BB0_73
; %bb.72:
	ds_read2_b32 v[4:5], v2 offset1:16
	s_mov_b64 s[2:3], exec
	s_waitcnt lgkmcnt(0)
	v_add_u32_e32 v4, vcc, v5, v4
	ds_write_b32 v2, v4
BB0_73:
	s_or_b64 exec, exec, s[4:5]
	s_waitcnt lgkmcnt(0)
	s_barrier
	s_waitcnt lgkmcnt(0)
	s_and_saveexec_b64 s[4:5], s[2:3]
	s_cbranch_execz BB0_75
; %bb.74:
	s_movk_i32 s2, 0x800
	v_add_u32_e32 v4, vcc, s2, v2
	ds_read2_b32 v[4:5], v4 offset1:16
	s_waitcnt lgkmcnt(0)
	v_add_u32_e32 v4, vcc, v5, v4
	ds_write_b32 v2, v4 offset:2048
BB0_75:
	s_or_b64 exec, exec, s[4:5]
	s_mov_b64 s[2:3], 0
	s_waitcnt lgkmcnt(0)
	s_barrier
	s_waitcnt lgkmcnt(0)
	s_and_saveexec_b64 s[4:5], s[0:1]
	s_cbranch_execz BB0_77
; %bb.76:
	s_movk_i32 s7, 0x1000
	v_add_u32_e32 v4, vcc, s7, v2
	ds_read2_b32 v[4:5], v4 offset1:16
	s_mov_b64 s[2:3], exec
	s_waitcnt lgkmcnt(0)
	v_add_u32_e32 v4, vcc, v5, v4
	ds_write_b32 v2, v4 offset:4096
BB0_77:
	s_or_b64 exec, exec, s[4:5]
	s_waitcnt lgkmcnt(0)
	s_barrier
	s_waitcnt lgkmcnt(0)
	s_and_saveexec_b64 s[4:5], s[2:3]
	s_cbranch_execz BB0_79
; %bb.78:
	s_movk_i32 s2, 0x1800
	v_add_u32_e32 v4, vcc, s2, v2
	ds_read2_b32 v[4:5], v4 offset1:16
	s_waitcnt lgkmcnt(0)
	v_add_u32_e32 v4, vcc, v5, v4
	ds_write_b32 v2, v4 offset:6144
BB0_79:
	s_or_b64 exec, exec, s[4:5]
	s_waitcnt lgkmcnt(0)
	s_barrier
	s_waitcnt lgkmcnt(0)
	s_and_saveexec_b64 s[2:3], s[0:1]
	s_cbranch_execz BB0_81
; %bb.80:
	s_movk_i32 s0, 0x2000
	v_add_u32_e32 v4, vcc, s0, v2
	ds_read2_b32 v[4:5], v4 offset1:16
	s_waitcnt lgkmcnt(0)
	v_add_u32_e32 v4, vcc, v5, v4
	ds_write_b32 v2, v4 offset:8192
BB0_81:
	s_or_b64 exec, exec, s[2:3]
	v_cmp_gt_u32_e64 s[0:1], 8, v3
	s_mov_b64 s[2:3], 0
	s_waitcnt lgkmcnt(0)
	s_barrier
	s_waitcnt lgkmcnt(0)
	s_and_saveexec_b64 s[4:5], s[0:1]
	s_cbranch_execz BB0_83
; %bb.82:
	ds_read2_b32 v[4:5], v2 offset1:8
	s_mov_b64 s[2:3], exec
	s_waitcnt lgkmcnt(0)
	v_add_u32_e32 v4, vcc, v5, v4
	ds_write_b32 v2, v4
BB0_83:
	s_or_b64 exec, exec, s[4:5]
	s_waitcnt lgkmcnt(0)
	s_barrier
	s_waitcnt lgkmcnt(0)
	s_and_saveexec_b64 s[4:5], s[2:3]
	s_cbranch_execz BB0_85
; %bb.84:
	s_movk_i32 s2, 0x800
	v_add_u32_e32 v4, vcc, s2, v2
	ds_read2_b32 v[4:5], v4 offset1:8
	s_waitcnt lgkmcnt(0)
	v_add_u32_e32 v4, vcc, v5, v4
	ds_write_b32 v2, v4 offset:2048
BB0_85:
	s_or_b64 exec, exec, s[4:5]
	s_mov_b64 s[2:3], 0
	s_waitcnt lgkmcnt(0)
	s_barrier
	s_waitcnt lgkmcnt(0)
	s_and_saveexec_b64 s[4:5], s[0:1]
	s_cbranch_execz BB0_87
; %bb.86:
	s_movk_i32 s7, 0x1000
	v_add_u32_e32 v4, vcc, s7, v2
	ds_read2_b32 v[4:5], v4 offset1:8
	s_mov_b64 s[2:3], exec
	s_waitcnt lgkmcnt(0)
	v_add_u32_e32 v4, vcc, v5, v4
	ds_write_b32 v2, v4 offset:4096
BB0_87:
	s_or_b64 exec, exec, s[4:5]
	s_waitcnt lgkmcnt(0)
	s_barrier
	s_waitcnt lgkmcnt(0)
	s_and_saveexec_b64 s[4:5], s[2:3]
	s_cbranch_execz BB0_89
; %bb.88:
	s_movk_i32 s2, 0x1800
	v_add_u32_e32 v4, vcc, s2, v2
	ds_read2_b32 v[4:5], v4 offset1:8
	s_waitcnt lgkmcnt(0)
	v_add_u32_e32 v4, vcc, v5, v4
	ds_write_b32 v2, v4 offset:6144
BB0_89:
	s_or_b64 exec, exec, s[4:5]
	s_waitcnt lgkmcnt(0)
	s_barrier
	s_waitcnt lgkmcnt(0)
	s_and_saveexec_b64 s[2:3], s[0:1]
	s_cbranch_execz BB0_91
; %bb.90:
	s_movk_i32 s0, 0x2000
	v_add_u32_e32 v4, vcc, s0, v2
	ds_read2_b32 v[4:5], v4 offset1:8
	s_waitcnt lgkmcnt(0)
	v_add_u32_e32 v4, vcc, v5, v4
	ds_write_b32 v2, v4 offset:8192
BB0_91:
	s_or_b64 exec, exec, s[2:3]
	v_cmp_gt_u32_e64 s[0:1], 4, v3
	s_mov_b64 s[2:3], 0
	s_waitcnt lgkmcnt(0)
	s_barrier
	s_waitcnt lgkmcnt(0)
	s_and_saveexec_b64 s[4:5], s[0:1]
	s_cbranch_execz BB0_93
; %bb.92:
	ds_read2_b32 v[4:5], v2 offset1:4
	s_mov_b64 s[2:3], exec
	s_waitcnt lgkmcnt(0)
	v_add_u32_e32 v4, vcc, v5, v4
	ds_write_b32 v2, v4
BB0_93:
	s_or_b64 exec, exec, s[4:5]
	s_waitcnt lgkmcnt(0)
	s_barrier
	s_waitcnt lgkmcnt(0)
	s_and_saveexec_b64 s[4:5], s[2:3]
	s_cbranch_execz BB0_95
; %bb.94:
	s_movk_i32 s2, 0x800
	v_add_u32_e32 v4, vcc, s2, v2
	ds_read2_b32 v[4:5], v4 offset1:4
	s_waitcnt lgkmcnt(0)
	v_add_u32_e32 v4, vcc, v5, v4
	ds_write_b32 v2, v4 offset:2048
BB0_95:
	s_or_b64 exec, exec, s[4:5]
	s_mov_b64 s[2:3], 0
	s_waitcnt lgkmcnt(0)
	s_barrier
	s_waitcnt lgkmcnt(0)
	s_and_saveexec_b64 s[4:5], s[0:1]
	s_cbranch_execz BB0_97
; %bb.96:
	s_movk_i32 s7, 0x1000
	v_add_u32_e32 v4, vcc, s7, v2
	ds_read2_b32 v[4:5], v4 offset1:4
	s_mov_b64 s[2:3], exec
	s_waitcnt lgkmcnt(0)
	v_add_u32_e32 v4, vcc, v5, v4
	ds_write_b32 v2, v4 offset:4096
BB0_97:
	s_or_b64 exec, exec, s[4:5]
	s_waitcnt lgkmcnt(0)
	s_barrier
	s_waitcnt lgkmcnt(0)
	s_and_saveexec_b64 s[4:5], s[2:3]
	s_cbranch_execz BB0_99
; %bb.98:
	s_movk_i32 s2, 0x1800
	v_add_u32_e32 v4, vcc, s2, v2
	ds_read2_b32 v[4:5], v4 offset1:4
	s_waitcnt lgkmcnt(0)
	v_add_u32_e32 v4, vcc, v5, v4
	ds_write_b32 v2, v4 offset:6144
BB0_99:
	s_or_b64 exec, exec, s[4:5]
	s_waitcnt lgkmcnt(0)
	s_barrier
	s_waitcnt lgkmcnt(0)
	s_and_saveexec_b64 s[2:3], s[0:1]
	s_cbranch_execz BB0_101
; %bb.100:
	s_movk_i32 s0, 0x2000
	v_add_u32_e32 v4, vcc, s0, v2
	ds_read2_b32 v[4:5], v4 offset1:4
	s_waitcnt lgkmcnt(0)
	v_add_u32_e32 v4, vcc, v5, v4
	ds_write_b32 v2, v4 offset:8192
BB0_101:
	s_or_b64 exec, exec, s[2:3]
	v_cmp_gt_u32_e64 s[0:1], 2, v3
	s_mov_b64 s[2:3], 0
	s_waitcnt lgkmcnt(0)
	s_barrier
	s_waitcnt lgkmcnt(0)
	s_and_saveexec_b64 s[4:5], s[0:1]
	s_cbranch_execz BB0_103
; %bb.102:
	ds_read2_b32 v[4:5], v2 offset1:2
	s_mov_b64 s[2:3], exec
	s_waitcnt lgkmcnt(0)
	v_add_u32_e32 v4, vcc, v5, v4
	ds_write_b32 v2, v4
BB0_103:
	s_or_b64 exec, exec, s[4:5]
	s_waitcnt lgkmcnt(0)
	s_barrier
	s_waitcnt lgkmcnt(0)
	s_and_saveexec_b64 s[4:5], s[2:3]
	s_cbranch_execz BB0_105
; %bb.104:
	s_movk_i32 s2, 0x800
	v_add_u32_e32 v4, vcc, s2, v2
	ds_read2_b32 v[4:5], v4 offset1:2
	s_waitcnt lgkmcnt(0)
	v_add_u32_e32 v4, vcc, v5, v4
	ds_write_b32 v2, v4 offset:2048
BB0_105:
	s_or_b64 exec, exec, s[4:5]
	s_mov_b64 s[2:3], 0
	s_waitcnt lgkmcnt(0)
	s_barrier
	s_waitcnt lgkmcnt(0)
	s_and_saveexec_b64 s[4:5], s[0:1]
	s_cbranch_execz BB0_107
; %bb.106:
	s_movk_i32 s7, 0x1000
	v_add_u32_e32 v4, vcc, s7, v2
	ds_read2_b32 v[4:5], v4 offset1:2
	s_mov_b64 s[2:3], exec
	s_waitcnt lgkmcnt(0)
	v_add_u32_e32 v4, vcc, v5, v4
	ds_write_b32 v2, v4 offset:4096
BB0_107:
	s_or_b64 exec, exec, s[4:5]
	s_waitcnt lgkmcnt(0)
	s_barrier
	s_waitcnt lgkmcnt(0)
	s_and_saveexec_b64 s[4:5], s[2:3]
	s_cbranch_execz BB0_109
; %bb.108:
	s_movk_i32 s2, 0x1800
	v_add_u32_e32 v4, vcc, s2, v2
	ds_read2_b32 v[4:5], v4 offset1:2
	s_waitcnt lgkmcnt(0)
	v_add_u32_e32 v4, vcc, v5, v4
	ds_write_b32 v2, v4 offset:6144
BB0_109:
	s_or_b64 exec, exec, s[4:5]
	s_waitcnt lgkmcnt(0)
	s_barrier
	s_waitcnt lgkmcnt(0)
	s_and_saveexec_b64 s[2:3], s[0:1]
	s_cbranch_execz BB0_111
; %bb.110:
	s_movk_i32 s0, 0x2000
	v_add_u32_e32 v4, vcc, s0, v2
	ds_read2_b32 v[4:5], v4 offset1:2
	s_waitcnt lgkmcnt(0)
	v_add_u32_e32 v4, vcc, v5, v4
	ds_write_b32 v2, v4 offset:8192
BB0_111:
	s_or_b64 exec, exec, s[2:3]
	v_cmp_eq_u32_e32 vcc, 0, v3
	s_mov_b64 s[2:3], 0
	s_waitcnt lgkmcnt(0)
	s_barrier
	s_waitcnt lgkmcnt(0)
	s_and_saveexec_b64 s[4:5], vcc
	s_cbranch_execz BB0_113
; %bb.112:
	ds_read2_b32 v[3:4], v2 offset1:1
	s_mov_b64 s[2:3], exec
	s_waitcnt lgkmcnt(0)
	v_add_u32_e64 v3, s[0:1], v4, v3
	ds_write_b32 v2, v3
BB0_113:
	s_or_b64 exec, exec, s[4:5]
	s_waitcnt lgkmcnt(0)
	s_barrier
	s_waitcnt lgkmcnt(0)
	s_and_saveexec_b64 s[4:5], s[2:3]
	s_cbranch_execz BB0_115
; %bb.114:
	v_or_b32_e32 v3, 0x800, v2
	ds_read2_b32 v[3:4], v3 offset1:1
	s_waitcnt lgkmcnt(0)
	v_add_u32_e64 v3, s[0:1], v4, v3
	ds_write_b32 v2, v3 offset:2048
BB0_115:
	s_or_b64 exec, exec, s[4:5]
	s_mov_b64 s[0:1], 0
	s_waitcnt lgkmcnt(0)
	s_barrier
	s_waitcnt lgkmcnt(0)
	s_and_saveexec_b64 s[2:3], vcc
	s_cbranch_execz BB0_117
; %bb.116:
	v_or_b32_e32 v3, 0x1000, v2
	ds_read2_b32 v[3:4], v3 offset1:1
	s_mov_b64 s[0:1], exec
	s_waitcnt lgkmcnt(0)
	v_add_u32_e32 v3, vcc, v4, v3
	ds_write_b32 v2, v3 offset:4096
BB0_117:
	s_or_b64 exec, exec, s[2:3]
	s_waitcnt lgkmcnt(0)
	s_barrier
	s_waitcnt lgkmcnt(0)
	s_and_saveexec_b64 s[2:3], s[0:1]
	s_cbranch_execz BB0_119
; %bb.118:
	v_or_b32_e32 v3, 0x1800, v2
	ds_read2_b32 v[3:4], v3 offset1:1
	s_waitcnt lgkmcnt(0)
	v_add_u32_e32 v3, vcc, v4, v3
	ds_write_b32 v2, v3 offset:6144
BB0_119:
	s_or_b64 exec, exec, s[2:3]
	v_and_b32_e32 v3, 63, v0
	v_cmp_eq_u32_e32 vcc, 0, v3
	s_waitcnt lgkmcnt(0)
	s_barrier
	s_waitcnt lgkmcnt(0)
	s_and_saveexec_b64 s[0:1], vcc
	s_cbranch_execz BB0_121
; %bb.120:
	v_or_b32_e32 v3, 0x2000, v2
	ds_read2_b32 v[3:4], v3 offset1:1
	s_waitcnt lgkmcnt(0)
	v_add_u32_e32 v3, vcc, v4, v3
	ds_write_b32 v2, v3 offset:8192
BB0_121:
	s_or_b64 exec, exec, s[0:1]
	v_cmp_gt_u32_e32 vcc, 20, v0
	s_waitcnt lgkmcnt(0)
	s_barrier
	s_waitcnt lgkmcnt(0)
	s_and_saveexec_b64 s[0:1], vcc
	s_xor_b64 s[0:1], exec, s[0:1]
	s_cbranch_execz BB0_123
; %bb.122:
	v_lshlrev_b32_e32 v2, 9, v0
	s_mul_i32 s0, s6, 20
	s_mov_b32 s1, 0
	s_lshl_b64 s[0:1], s[0:1], 3
	ds_read_b32 v2, v2
	s_add_u32 s0, s8, s0
	v_lshlrev_b64 v[0:1], 3, v[0:1]
	s_addc_u32 s1, s9, s1
	v_mov_b32_e32 v3, s1
	v_add_u32_e32 v0, vcc, s0, v0
	v_addc_u32_e32 v1, vcc, v3, v1, vcc
	v_mov_b32_e32 v3, 0
	s_waitcnt lgkmcnt(0)
	flat_store_dwordx2 v[0:1], v[2:3]
BB0_123:
	s_endpgm
.Lfunc_end0:
	.size	gen_hists, .Lfunc_end0-gen_hists
                                        ; -- End function
	.section	.AMDGPU.csdata
; Kernel info:
; codeLenInByte = 5108
; NumSgprs: 42
; NumVgprs: 15
; ScratchSize: 0
; MemoryBound: 0
; FloatMode: 192
; IeeeMode: 1
; LDSByteSize: 10240 bytes/workgroup (compile time only)
; SGPRBlocks: 5
; VGPRBlocks: 3
; NumSGPRsForWavesPerEU: 42
; NumVGPRsForWavesPerEU: 15
; Occupancy: 10
; WaveLimiterHint : 1
; COMPUTE_PGM_RSRC2:USER_SGPR: 6
; COMPUTE_PGM_RSRC2:TRAP_HANDLER: 0
; COMPUTE_PGM_RSRC2:TGID_X_EN: 1
; COMPUTE_PGM_RSRC2:TGID_Y_EN: 0
; COMPUTE_PGM_RSRC2:TGID_Z_EN: 0
; COMPUTE_PGM_RSRC2:TIDIG_COMP_CNT: 0
	.ident	"clang version 11.0.0 (/src/external/llvm-project/clang 0383ad1cfb0a8e05b0a020e8632400194628b243)"
	.section	".note.GNU-stack"
	.addrsig
	.addrsig_sym gen_hists.warp_hists
	.amd_amdgpu_isa "amdgcn-amd-amdhsa--gfx803"
	.amd_amdgpu_hsa_metadata
---
Version:         [ 1, 0 ]
Kernels:
  - Name:            gen_hists
    SymbolName:      'gen_hists@kd'
    Language:        OpenCL C
    LanguageVersion: [ 1, 2 ]
    Args:
      - Name:            histograms
        TypeName:        'hist_t*'
        Size:            8
        Align:           8
        ValueKind:       GlobalBuffer
        ValueType:       U64
        AddrSpaceQual:   Global
        AccQual:         Default
      - Name:            all_x_data
        TypeName:        'float*'
        Size:            8
        Align:           8
        ValueKind:       GlobalBuffer
        ValueType:       F32
        AddrSpaceQual:   Global
        AccQual:         Default
      - Name:            dev_binb
        TypeName:        'float*'
        Size:            8
        Align:           8
        ValueKind:       GlobalBuffer
        ValueType:       F32
        AddrSpaceQual:   Constant
        AccQual:         Default
        IsConst:         true
      - Name:            NUM_SETS
        TypeName:        int
        Size:            4
        Align:           4
        ValueKind:       ByValue
        ValueType:       I32
        AccQual:         Default
      - Name:            NUM_ELEMENTS
        TypeName:        int
        Size:            4
        Align:           4
        ValueKind:       ByValue
        ValueType:       I32
        AccQual:         Default
      - Size:            8
        Align:           8
        ValueKind:       HiddenGlobalOffsetX
        ValueType:       I64
      - Size:            8
        Align:           8
        ValueKind:       HiddenGlobalOffsetY
        ValueType:       I64
      - Size:            8
        Align:           8
        ValueKind:       HiddenGlobalOffsetZ
        ValueType:       I64
      - Size:            8
        Align:           8
        ValueKind:       HiddenNone
        ValueType:       I8
        AddrSpaceQual:   Global
      - Size:            8
        Align:           8
        ValueKind:       HiddenNone
        ValueType:       I8
        AddrSpaceQual:   Global
      - Size:            8
        Align:           8
        ValueKind:       HiddenNone
        ValueType:       I8
        AddrSpaceQual:   Global
      - Size:            8
        Align:           8
        ValueKind:       HiddenMultiGridSyncArg
        ValueType:       I8
        AddrSpaceQual:   Global
    CodeProps:
      KernargSegmentSize: 88
      GroupSegmentFixedSize: 10240
      PrivateSegmentFixedSize: 0
      KernargSegmentAlign: 8
      WavefrontSize:   64
      NumSGPRs:        42
      NumVGPRs:        15
      MaxFlatWorkGroupSize: 256
...

	.end_amd_amdgpu_hsa_metadata
