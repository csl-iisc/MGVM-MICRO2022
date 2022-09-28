	.text
	.hsa_code_object_version 2,1
	.hsa_code_object_isa 8,0,3,"AMD","AMDGPU"
	.protected	transpose_tensor ; -- Begin function transpose_tensor
	.globl	transpose_tensor
	.p2align	8
	.type	transpose_tensor,@function
	.amdgpu_hsa_kernel transpose_tensor
transpose_tensor:                       ; @transpose_tensor
transpose_tensor$local:
	.amd_kernel_code_t
		amd_code_version_major = 1
		amd_code_version_minor = 2
		amd_machine_kind = 1
		amd_machine_version_major = 8
		amd_machine_version_minor = 0
		amd_machine_version_stepping = 3
		kernel_code_entry_byte_offset = 256
		kernel_code_prefetch_byte_size = 0
		granulated_workitem_vgpr_count = 4
		granulated_wavefront_sgpr_count = 3
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
		user_sgpr_count = 8
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
		enable_sgpr_dispatch_ptr = 1
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
		workgroup_group_segment_byte_size = 0
		gds_segment_byte_size = 0
		kernarg_segment_byte_size = 120
		workgroup_fbarrier_count = 0
		wavefront_sgpr_count = 26
		workitem_vgpr_count = 18
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
	s_load_dwordx8 s[12:19], s[6:7], 0x0
	s_load_dword s0, s[4:5], 0x4
	s_waitcnt lgkmcnt(0)
	s_and_b32 s0, s0, 0xffff
	s_mul_i32 s8, s8, s0
	s_load_dword s2, s[6:7], 0x38
	s_load_dwordx2 s[0:1], s[6:7], 0x40
	v_add_u32_e32 v0, vcc, s8, v0
	s_waitcnt lgkmcnt(0)
	s_cmp_lt_i32 s2, 1
	v_add_u32_e32 v1, vcc, s0, v0
	s_mov_b32 s0, 1
	s_cbranch_scc1 BB0_10
; %bb.1:                                ; %.preheader3.preheader
	v_mul_lo_u32 v2, v1, s2
	s_load_dwordx4 s[8:11], s[6:7], 0x20
	s_load_dwordx4 s[20:23], s[6:7], 0x30
	s_mov_b64 s[4:5], s[18:19]
	s_mov_b64 s[6:7], s[18:19]
	v_ashrrev_i32_e32 v3, 31, v2
	v_lshlrev_b64 v[3:4], 2, v[2:3]
	s_waitcnt lgkmcnt(0)
	v_mov_b32_e32 v2, s11
	v_add_u32_e32 v0, vcc, s10, v3
	v_addc_u32_e32 v2, vcc, v2, v4, vcc
	v_mov_b32_e32 v5, s21
	v_add_u32_e32 v3, vcc, s20, v3
	v_addc_u32_e32 v4, vcc, v5, v4, vcc
	s_mov_b32 s1, s2
BB0_2:                                  ; %.preheader3
                                        ; =>This Inner Loop Header: Depth=1
	s_load_dword s3, s[6:7], 0x0
	s_add_i32 s1, s1, -1
	s_add_u32 s6, s6, 4
	s_addc_u32 s7, s7, 0
	s_cmp_lg_u32 s1, 0
	s_waitcnt lgkmcnt(0)
	s_mul_i32 s0, s3, s0
	s_cbranch_scc1 BB0_2
; %bb.3:                                ; %.preheader2.preheader
	v_mov_b32_e32 v6, v4
	v_mov_b32_e32 v8, s0
	s_mov_b32 s3, 0x4f800000
	v_mov_b32_e32 v5, v3
	s_mov_b32 s6, s2
	v_mov_b32_e32 v7, v1
BB0_4:                                  ; %.preheader2
                                        ; =>This Inner Loop Header: Depth=1
	v_mov_b32_e32 v10, s5
	v_mov_b32_e32 v9, s4
	flat_load_dword v9, v[9:10]
	v_ashrrev_i32_e32 v10, 31, v8
	v_add_u32_e32 v8, vcc, v10, v8
	v_xor_b32_e32 v8, v8, v10
	v_ashrrev_i32_e32 v11, 31, v7
	v_add_u32_e32 v12, vcc, v7, v11
	v_xor_b32_e32 v12, v12, v11
	s_add_i32 s6, s6, -1
	s_waitcnt vmcnt(0) lgkmcnt(0)
	v_ashrrev_i32_e32 v13, 31, v9
	v_add_u32_e32 v9, vcc, v9, v13
	v_xor_b32_e32 v9, v9, v13
	v_xor_b32_e32 v10, v10, v13
	v_cvt_f32_u32_e32 v13, v9
	v_rcp_iflag_f32_e32 v13, v13
	v_mul_f32_e32 v13, s3, v13
	v_cvt_u32_f32_e32 v13, v13
	v_mul_lo_u32 v14, v13, v9
	v_mul_hi_u32 v15, v13, v9
	v_sub_u32_e32 v16, vcc, 0, v14
	v_cmp_eq_u32_e64 s[0:1], 0, v15
	v_cndmask_b32_e64 v14, v14, v16, s[0:1]
	v_mul_hi_u32 v14, v14, v13
	v_add_u32_e32 v15, vcc, v14, v13
	v_subrev_u32_e32 v13, vcc, v14, v13
	v_cndmask_b32_e64 v13, v13, v15, s[0:1]
	v_mul_hi_u32 v13, v13, v8
	v_mul_lo_u32 v14, v13, v9
	v_add_u32_e32 v15, vcc, 1, v13
	v_add_u32_e32 v16, vcc, -1, v13
	v_sub_u32_e32 v17, vcc, v8, v14
	v_cmp_ge_u32_e32 vcc, v8, v14
	v_cmp_ge_u32_e64 s[0:1], v17, v9
	s_and_b64 s[0:1], s[0:1], vcc
	v_cndmask_b32_e64 v8, v13, v15, s[0:1]
	v_cndmask_b32_e32 v8, v16, v8, vcc
	v_xor_b32_e32 v8, v8, v10
	v_sub_u32_e32 v8, vcc, v8, v10
	v_ashrrev_i32_e32 v9, 31, v8
	v_xor_b32_e32 v10, v11, v9
	v_add_u32_e32 v11, vcc, v9, v8
	v_xor_b32_e32 v9, v11, v9
	v_cvt_f32_u32_e32 v11, v9
	v_rcp_iflag_f32_e32 v11, v11
	v_mul_f32_e32 v11, s3, v11
	v_cvt_u32_f32_e32 v11, v11
	v_mul_lo_u32 v13, v11, v9
	v_mul_hi_u32 v14, v11, v9
	v_sub_u32_e32 v15, vcc, 0, v13
	v_cmp_eq_u32_e64 s[0:1], 0, v14
	v_cndmask_b32_e64 v13, v13, v15, s[0:1]
	v_mul_hi_u32 v13, v13, v11
	v_add_u32_e32 v14, vcc, v13, v11
	v_subrev_u32_e32 v11, vcc, v13, v11
	v_cndmask_b32_e64 v11, v11, v14, s[0:1]
	v_mul_hi_u32 v11, v11, v12
	v_mul_lo_u32 v13, v11, v9
	v_add_u32_e32 v14, vcc, 1, v11
	v_add_u32_e32 v15, vcc, -1, v11
	v_sub_u32_e32 v16, vcc, v12, v13
	v_cmp_ge_u32_e32 vcc, v12, v13
	v_cmp_ge_u32_e64 s[0:1], v16, v9
	s_and_b64 s[0:1], s[0:1], vcc
	v_cndmask_b32_e64 v9, v11, v14, s[0:1]
	v_cndmask_b32_e32 v9, v15, v9, vcc
	v_xor_b32_e32 v9, v9, v10
	v_sub_u32_e32 v9, vcc, v9, v10
	v_mul_lo_u32 v10, v9, v8
	s_add_u32 s4, s4, 4
	s_addc_u32 s5, s5, 0
	flat_store_dword v[5:6], v9
	v_add_u32_e32 v5, vcc, 4, v5
	v_addc_u32_e32 v6, vcc, 0, v6, vcc
	s_cmp_lg_u32 s6, 0
	v_sub_u32_e32 v7, vcc, v7, v10
	s_cbranch_scc1 BB0_4
; %bb.5:                                ; %.preheader1.preheader
	s_mov_b32 s0, s2
BB0_6:                                  ; %.preheader1
                                        ; =>This Inner Loop Header: Depth=1
	v_mov_b32_e32 v5, s8
	v_mov_b32_e32 v6, s9
	flat_load_dword v5, v[5:6]
	flat_load_dword v7, v[3:4]
	v_add_u32_e32 v3, vcc, 4, v3
	s_add_i32 s0, s0, -1
	s_add_u32 s8, s8, 4
	v_addc_u32_e32 v4, vcc, 0, v4, vcc
	s_addc_u32 s9, s9, 0
	s_cmp_lg_u32 s0, 0
	s_waitcnt vmcnt(1) lgkmcnt(1)
	v_ashrrev_i32_e32 v6, 31, v5
	v_lshlrev_b64 v[5:6], 2, v[5:6]
	v_add_u32_e32 v5, vcc, v0, v5
	v_addc_u32_e32 v6, vcc, v2, v6, vcc
	s_waitcnt vmcnt(0) lgkmcnt(0)
	flat_store_dword v[5:6], v7
	s_cbranch_scc1 BB0_6
; %bb.7:                                ; %.preheader.preheader
	s_add_i32 s0, s2, -1
	v_mov_b32_e32 v3, 0
	v_mov_b32_e32 v4, 1
BB0_8:                                  ; %.preheader
                                        ; =>This Inner Loop Header: Depth=1
	s_ashr_i32 s1, s0, 31
	s_lshl_b64 s[2:3], s[0:1], 2
	s_add_u32 s4, s16, s2
	v_mov_b32_e32 v6, s3
	v_add_u32_e32 v5, vcc, s2, v0
	v_addc_u32_e32 v6, vcc, v2, v6, vcc
	s_addc_u32 s5, s17, s3
	flat_load_dword v7, v[5:6]
	v_mov_b32_e32 v6, s5
	v_mov_b32_e32 v5, s4
	flat_load_dword v5, v[5:6]
	s_add_i32 s0, s0, -1
	s_cmp_lg_u32 s0, -1
	s_waitcnt vmcnt(1) lgkmcnt(1)
	v_mul_lo_u32 v6, v7, v4
	s_waitcnt vmcnt(0) lgkmcnt(0)
	v_mul_lo_u32 v4, v5, v4
	v_add_u32_e32 v3, vcc, v6, v3
	s_cbranch_scc1 BB0_8
; %bb.9:
	v_ashrrev_i32_e32 v4, 31, v3
	s_branch BB0_11
BB0_10:
	v_mov_b32_e32 v3, 0
	v_mov_b32_e32 v4, 0
BB0_11:
	v_lshlrev_b64 v[2:3], 2, v[3:4]
	v_mov_b32_e32 v0, s13
	v_add_u32_e32 v2, vcc, s12, v2
	v_addc_u32_e32 v3, vcc, v0, v3, vcc
	flat_load_dword v2, v[2:3]
	v_mov_b32_e32 v0, 0
	v_ashrrev_i64 v[0:1], 30, v[0:1]
	v_mov_b32_e32 v3, s15
	v_add_u32_e32 v0, vcc, s14, v0
	v_addc_u32_e32 v1, vcc, v3, v1, vcc
	s_waitcnt vmcnt(0) lgkmcnt(0)
	flat_store_dword v[0:1], v2
	s_endpgm
.Lfunc_end0:
	.size	transpose_tensor, .Lfunc_end0-transpose_tensor
                                        ; -- End function
	.section	.AMDGPU.csdata
; Kernel info:
; codeLenInByte = 896
; NumSgprs: 26
; NumVgprs: 18
; ScratchSize: 0
; MemoryBound: 0
; FloatMode: 192
; IeeeMode: 1
; LDSByteSize: 0 bytes/workgroup (compile time only)
; SGPRBlocks: 3
; VGPRBlocks: 4
; NumSGPRsForWavesPerEU: 26
; NumVGPRsForWavesPerEU: 18
; Occupancy: 10
; WaveLimiterHint : 1
; COMPUTE_PGM_RSRC2:USER_SGPR: 8
; COMPUTE_PGM_RSRC2:TRAP_HANDLER: 0
; COMPUTE_PGM_RSRC2:TGID_X_EN: 1
; COMPUTE_PGM_RSRC2:TGID_Y_EN: 0
; COMPUTE_PGM_RSRC2:TGID_Z_EN: 0
; COMPUTE_PGM_RSRC2:TIDIG_COMP_CNT: 0
	.text
	.protected	rotate_tensor   ; -- Begin function rotate_tensor
	.globl	rotate_tensor
	.p2align	8
	.type	rotate_tensor,@function
	.amdgpu_hsa_kernel rotate_tensor
rotate_tensor:                          ; @rotate_tensor
rotate_tensor$local:
	.amd_kernel_code_t
		amd_code_version_major = 1
		amd_code_version_minor = 2
		amd_machine_kind = 1
		amd_machine_version_major = 8
		amd_machine_version_minor = 0
		amd_machine_version_stepping = 3
		kernel_code_entry_byte_offset = 256
		kernel_code_prefetch_byte_size = 0
		granulated_workitem_vgpr_count = 4
		granulated_wavefront_sgpr_count = 2
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
		user_sgpr_count = 8
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
		enable_sgpr_dispatch_ptr = 1
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
		workgroup_group_segment_byte_size = 0
		gds_segment_byte_size = 0
		kernarg_segment_byte_size = 112
		workgroup_fbarrier_count = 0
		wavefront_sgpr_count = 18
		workitem_vgpr_count = 18
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
	s_load_dword s2, s[6:7], 0x30
	s_load_dwordx2 s[0:1], s[6:7], 0x38
	s_waitcnt lgkmcnt(0)
	s_load_dword s1, s[4:5], 0x4
	s_waitcnt lgkmcnt(0)
	s_and_b32 s1, s1, 0xffff
	s_mul_i32 s8, s8, s1
	v_add_u32_e32 v0, vcc, s8, v0
	v_add_u32_e32 v0, vcc, s0, v0
	v_mul_lo_u32 v1, v0, s2
	s_load_dwordx4 s[12:15], s[6:7], 0x20
	s_load_dwordx8 s[4:11], s[6:7], 0x0
	s_cmp_lt_i32 s2, 1
	s_mov_b32 s0, 1
	v_ashrrev_i32_e32 v2, 31, v1
	v_lshlrev_b64 v[3:4], 2, v[1:2]
	s_waitcnt lgkmcnt(0)
	v_mov_b32_e32 v2, s13
	v_add_u32_e32 v1, vcc, s12, v3
	v_addc_u32_e32 v2, vcc, v2, v4, vcc
	v_mov_b32_e32 v6, s15
	v_add_u32_e32 v3, vcc, s14, v3
	v_addc_u32_e32 v4, vcc, v6, v4, vcc
	s_cbranch_scc1 BB1_7
; %bb.1:                                ; %.preheader3.preheader
	v_mov_b32_e32 v5, v0
	s_mov_b64 s[12:13], s[10:11]
	s_mov_b32 s1, s2
BB1_2:                                  ; %.preheader3
                                        ; =>This Inner Loop Header: Depth=1
	s_load_dword s3, s[12:13], 0x0
	s_add_i32 s1, s1, -1
	s_add_u32 s12, s12, 4
	s_addc_u32 s13, s13, 0
	s_cmp_lg_u32 s1, 0
	s_waitcnt lgkmcnt(0)
	s_mul_i32 s0, s3, s0
	s_cbranch_scc1 BB1_2
; %bb.3:                                ; %.preheader2.preheader
	v_mov_b32_e32 v7, v4
	v_mov_b32_e32 v8, s0
	s_mov_b32 s3, 0x4f800000
	v_mov_b32_e32 v6, v3
	s_mov_b32 s12, s2
BB1_4:                                  ; %.preheader2
                                        ; =>This Inner Loop Header: Depth=1
	v_mov_b32_e32 v9, s10
	v_mov_b32_e32 v10, s11
	flat_load_dword v9, v[9:10]
	v_ashrrev_i32_e32 v10, 31, v8
	v_add_u32_e32 v8, vcc, v10, v8
	v_xor_b32_e32 v8, v8, v10
	v_ashrrev_i32_e32 v11, 31, v5
	v_add_u32_e32 v12, vcc, v5, v11
	v_xor_b32_e32 v12, v12, v11
	s_add_i32 s12, s12, -1
	s_waitcnt vmcnt(0) lgkmcnt(0)
	v_ashrrev_i32_e32 v13, 31, v9
	v_add_u32_e32 v9, vcc, v9, v13
	v_xor_b32_e32 v9, v9, v13
	v_xor_b32_e32 v10, v10, v13
	v_cvt_f32_u32_e32 v13, v9
	v_rcp_iflag_f32_e32 v13, v13
	v_mul_f32_e32 v13, s3, v13
	v_cvt_u32_f32_e32 v13, v13
	v_mul_lo_u32 v14, v13, v9
	v_mul_hi_u32 v15, v13, v9
	v_sub_u32_e32 v16, vcc, 0, v14
	v_cmp_eq_u32_e64 s[0:1], 0, v15
	v_cndmask_b32_e64 v14, v14, v16, s[0:1]
	v_mul_hi_u32 v14, v14, v13
	v_add_u32_e32 v15, vcc, v14, v13
	v_subrev_u32_e32 v13, vcc, v14, v13
	v_cndmask_b32_e64 v13, v13, v15, s[0:1]
	v_mul_hi_u32 v13, v13, v8
	v_mul_lo_u32 v14, v13, v9
	v_add_u32_e32 v15, vcc, 1, v13
	v_add_u32_e32 v16, vcc, -1, v13
	v_sub_u32_e32 v17, vcc, v8, v14
	v_cmp_ge_u32_e32 vcc, v8, v14
	v_cmp_ge_u32_e64 s[0:1], v17, v9
	s_and_b64 s[0:1], s[0:1], vcc
	v_cndmask_b32_e64 v8, v13, v15, s[0:1]
	v_cndmask_b32_e32 v8, v16, v8, vcc
	v_xor_b32_e32 v8, v8, v10
	v_sub_u32_e32 v8, vcc, v8, v10
	v_ashrrev_i32_e32 v9, 31, v8
	v_xor_b32_e32 v10, v11, v9
	v_add_u32_e32 v11, vcc, v9, v8
	v_xor_b32_e32 v9, v11, v9
	v_cvt_f32_u32_e32 v11, v9
	v_rcp_iflag_f32_e32 v11, v11
	v_mul_f32_e32 v11, s3, v11
	v_cvt_u32_f32_e32 v11, v11
	v_mul_lo_u32 v13, v11, v9
	v_mul_hi_u32 v14, v11, v9
	v_sub_u32_e32 v15, vcc, 0, v13
	v_cmp_eq_u32_e64 s[0:1], 0, v14
	v_cndmask_b32_e64 v13, v13, v15, s[0:1]
	v_mul_hi_u32 v13, v13, v11
	v_add_u32_e32 v14, vcc, v13, v11
	v_subrev_u32_e32 v11, vcc, v13, v11
	v_cndmask_b32_e64 v11, v11, v14, s[0:1]
	v_mul_hi_u32 v11, v11, v12
	v_mul_lo_u32 v13, v11, v9
	v_add_u32_e32 v14, vcc, 1, v11
	v_add_u32_e32 v15, vcc, -1, v11
	v_sub_u32_e32 v16, vcc, v12, v13
	v_cmp_ge_u32_e32 vcc, v12, v13
	v_cmp_ge_u32_e64 s[0:1], v16, v9
	s_and_b64 s[0:1], s[0:1], vcc
	v_cndmask_b32_e64 v9, v11, v14, s[0:1]
	v_cndmask_b32_e32 v9, v15, v9, vcc
	v_xor_b32_e32 v9, v9, v10
	v_sub_u32_e32 v9, vcc, v9, v10
	v_mul_lo_u32 v10, v9, v8
	s_add_u32 s10, s10, 4
	s_addc_u32 s11, s11, 0
	flat_store_dword v[6:7], v9
	v_add_u32_e32 v6, vcc, 4, v6
	v_addc_u32_e32 v7, vcc, 0, v7, vcc
	s_cmp_lg_u32 s12, 0
	v_sub_u32_e32 v5, vcc, v5, v10
	s_cbranch_scc1 BB1_4
; %bb.5:                                ; %.preheader1.preheader
	v_mov_b32_e32 v6, v4
	v_mov_b32_e32 v8, v2
	v_mov_b32_e32 v5, v3
	v_mov_b32_e32 v7, v1
	s_mov_b32 s0, s2
BB1_6:                                  ; %.preheader1
                                        ; =>This Inner Loop Header: Depth=1
	flat_load_dword v9, v[5:6]
	v_add_u32_e32 v5, vcc, 4, v5
	v_addc_u32_e32 v6, vcc, 0, v6, vcc
	s_add_i32 s0, s0, -1
	s_cmp_eq_u32 s0, 0
	s_waitcnt vmcnt(0) lgkmcnt(0)
	flat_store_dword v[7:8], v9
	v_add_u32_e32 v7, vcc, 4, v7
	v_addc_u32_e32 v8, vcc, 0, v8, vcc
	s_cbranch_scc0 BB1_6
BB1_7:                                  ; %Flow46
	s_add_i32 s0, s2, -1
	s_ashr_i32 s1, s0, 31
	s_lshl_b64 s[12:13], s[0:1], 2
	s_add_u32 s14, s8, s12
	s_addc_u32 s15, s9, s13
	v_mov_b32_e32 v5, s14
	v_mov_b32_e32 v6, s15
	flat_load_dword v11, v[5:6]
	v_mov_b32_e32 v6, s13
	v_add_u32_e32 v5, vcc, s12, v3
	v_addc_u32_e32 v6, vcc, v4, v6, vcc
	flat_load_dword v12, v[5:6]
	v_cmp_gt_i32_e64 s[10:11], s2, 0
	s_add_i32 s2, s2, -2
	s_ashr_i32 s3, s2, 31
	v_mov_b32_e32 v6, s13
	v_add_u32_e32 v5, vcc, s12, v1
	s_lshl_b64 s[2:3], s[2:3], 2
	v_addc_u32_e32 v6, vcc, v2, v6, vcc
	v_mov_b32_e32 v7, s3
	v_add_u32_e32 v3, vcc, s2, v3
	v_addc_u32_e32 v4, vcc, v4, v7, vcc
	v_mov_b32_e32 v8, s3
	v_add_u32_e32 v7, vcc, s2, v1
	s_add_u32 s12, s8, s2
	v_addc_u32_e32 v8, vcc, v2, v8, vcc
	s_addc_u32 s13, s9, s3
	v_mov_b32_e32 v9, s12
	v_mov_b32_e32 v10, s13
	s_waitcnt vmcnt(0) lgkmcnt(0)
	v_not_b32_e32 v12, v12
	v_add_u32_e32 v11, vcc, v11, v12
	flat_store_dword v[5:6], v11
	flat_load_dword v3, v[3:4]
	flat_load_dword v4, v[9:10]
	s_waitcnt vmcnt(1) lgkmcnt(1)
	v_not_b32_e32 v3, v3
	s_waitcnt vmcnt(0) lgkmcnt(0)
	v_add_u32_e32 v3, vcc, v4, v3
	s_andn2_b64 vcc, exec, s[10:11]
	flat_store_dword v[7:8], v3
	s_cbranch_vccnz BB1_11
; %bb.8:                                ; %.preheader.preheader
	v_mov_b32_e32 v3, 0
	v_mov_b32_e32 v4, 1
BB1_9:                                  ; %.preheader
                                        ; =>This Inner Loop Header: Depth=1
	s_ashr_i32 s1, s0, 31
	s_lshl_b64 s[2:3], s[0:1], 2
	v_mov_b32_e32 v6, s3
	v_add_u32_e32 v5, vcc, s2, v1
	s_add_u32 s10, s8, s2
	v_addc_u32_e32 v6, vcc, v2, v6, vcc
	s_addc_u32 s11, s9, s3
	flat_load_dword v7, v[5:6]
	v_mov_b32_e32 v5, s10
	v_mov_b32_e32 v6, s11
	flat_load_dword v5, v[5:6]
	s_add_i32 s0, s0, -1
	s_cmp_lg_u32 s0, -1
	s_waitcnt vmcnt(1) lgkmcnt(1)
	v_mul_lo_u32 v6, v7, v4
	s_waitcnt vmcnt(0) lgkmcnt(0)
	v_mul_lo_u32 v4, v5, v4
	v_add_u32_e32 v3, vcc, v6, v3
	s_cbranch_scc1 BB1_9
; %bb.10:
	v_ashrrev_i32_e32 v4, 31, v3
	s_branch BB1_12
BB1_11:
	v_mov_b32_e32 v3, 0
	v_mov_b32_e32 v4, 0
BB1_12:
	v_lshlrev_b64 v[1:2], 2, v[3:4]
	v_mov_b32_e32 v3, s5
	v_add_u32_e32 v1, vcc, s4, v1
	v_addc_u32_e32 v2, vcc, v3, v2, vcc
	flat_load_dword v3, v[1:2]
	v_mov_b32_e32 v1, 0
	v_mov_b32_e32 v2, v0
	v_ashrrev_i64 v[0:1], 30, v[1:2]
	v_mov_b32_e32 v2, s7
	v_add_u32_e32 v0, vcc, s6, v0
	v_addc_u32_e32 v1, vcc, v2, v1, vcc
	s_waitcnt vmcnt(0) lgkmcnt(0)
	flat_store_dword v[0:1], v3
	s_endpgm
.Lfunc_end1:
	.size	rotate_tensor, .Lfunc_end1-rotate_tensor
                                        ; -- End function
	.section	.AMDGPU.csdata
; Kernel info:
; codeLenInByte = 1056
; NumSgprs: 18
; NumVgprs: 18
; ScratchSize: 0
; MemoryBound: 0
; FloatMode: 192
; IeeeMode: 1
; LDSByteSize: 0 bytes/workgroup (compile time only)
; SGPRBlocks: 2
; VGPRBlocks: 4
; NumSGPRsForWavesPerEU: 18
; NumVGPRsForWavesPerEU: 18
; Occupancy: 10
; WaveLimiterHint : 1
; COMPUTE_PGM_RSRC2:USER_SGPR: 8
; COMPUTE_PGM_RSRC2:TRAP_HANDLER: 0
; COMPUTE_PGM_RSRC2:TGID_X_EN: 1
; COMPUTE_PGM_RSRC2:TGID_Y_EN: 0
; COMPUTE_PGM_RSRC2:TGID_Z_EN: 0
; COMPUTE_PGM_RSRC2:TIDIG_COMP_CNT: 0
	.text
	.protected	dilate_tensor   ; -- Begin function dilate_tensor
	.globl	dilate_tensor
	.p2align	8
	.type	dilate_tensor,@function
	.amdgpu_hsa_kernel dilate_tensor
dilate_tensor:                          ; @dilate_tensor
dilate_tensor$local:
	.amd_kernel_code_t
		amd_code_version_major = 1
		amd_code_version_minor = 2
		amd_machine_kind = 1
		amd_machine_version_major = 8
		amd_machine_version_minor = 0
		amd_machine_version_stepping = 3
		kernel_code_entry_byte_offset = 256
		kernel_code_prefetch_byte_size = 0
		granulated_workitem_vgpr_count = 4
		granulated_wavefront_sgpr_count = 3
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
		user_sgpr_count = 8
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
		enable_sgpr_dispatch_ptr = 1
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
		workgroup_group_segment_byte_size = 0
		gds_segment_byte_size = 0
		kernarg_segment_byte_size = 120
		workgroup_fbarrier_count = 0
		wavefront_sgpr_count = 30
		workitem_vgpr_count = 18
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
	s_load_dwordx8 s[12:19], s[6:7], 0x0
	s_load_dword s0, s[4:5], 0x4
	s_waitcnt lgkmcnt(0)
	s_and_b32 s0, s0, 0xffff
	s_mul_i32 s8, s8, s0
	s_load_dword s4, s[6:7], 0x38
	s_load_dwordx2 s[0:1], s[6:7], 0x40
	v_add_u32_e32 v0, vcc, s8, v0
	s_waitcnt lgkmcnt(0)
	v_cmp_gt_i32_e64 s[20:21], s4, 0
	v_add_u32_e32 v0, vcc, s0, v0
	v_mul_lo_u32 v1, v0, s4
	s_load_dwordx4 s[8:11], s[6:7], 0x20
	s_load_dwordx4 s[0:3], s[6:7], 0x30
	s_cmp_lt_i32 s4, 1
	v_ashrrev_i32_e32 v2, 31, v1
	v_lshlrev_b64 v[3:4], 2, v[1:2]
	s_waitcnt lgkmcnt(0)
	v_mov_b32_e32 v6, s1
	v_add_u32_e32 v3, vcc, s0, v3
	v_addc_u32_e32 v4, vcc, v6, v4, vcc
	s_mov_b32 s0, 1
	s_cbranch_scc1 BB2_5
; %bb.1:                                ; %.preheader3.preheader
	v_mov_b32_e32 v5, v0
	s_mov_b64 s[2:3], s[18:19]
	s_mov_b64 s[6:7], s[18:19]
	s_mov_b32 s1, s4
BB2_2:                                  ; %.preheader3
                                        ; =>This Inner Loop Header: Depth=1
	s_load_dword s5, s[6:7], 0x0
	s_add_i32 s1, s1, -1
	s_add_u32 s6, s6, 4
	s_addc_u32 s7, s7, 0
	s_cmp_lg_u32 s1, 0
	s_waitcnt lgkmcnt(0)
	s_mul_i32 s0, s5, s0
	s_cbranch_scc1 BB2_2
; %bb.3:                                ; %.preheader2.preheader
	v_mov_b32_e32 v7, v4
	v_mov_b32_e32 v8, s0
	s_mov_b32 s5, 0x4f800000
	v_mov_b32_e32 v6, v3
	s_mov_b32 s6, s4
BB2_4:                                  ; %.preheader2
                                        ; =>This Inner Loop Header: Depth=1
	v_mov_b32_e32 v10, s3
	v_mov_b32_e32 v9, s2
	flat_load_dword v9, v[9:10]
	v_ashrrev_i32_e32 v10, 31, v8
	v_add_u32_e32 v8, vcc, v10, v8
	v_xor_b32_e32 v8, v8, v10
	v_ashrrev_i32_e32 v11, 31, v5
	v_add_u32_e32 v12, vcc, v5, v11
	v_xor_b32_e32 v12, v12, v11
	s_add_i32 s6, s6, -1
	s_waitcnt vmcnt(0) lgkmcnt(0)
	v_ashrrev_i32_e32 v13, 31, v9
	v_add_u32_e32 v9, vcc, v9, v13
	v_xor_b32_e32 v9, v9, v13
	v_xor_b32_e32 v10, v10, v13
	v_cvt_f32_u32_e32 v13, v9
	v_rcp_iflag_f32_e32 v13, v13
	v_mul_f32_e32 v13, s5, v13
	v_cvt_u32_f32_e32 v13, v13
	v_mul_lo_u32 v14, v13, v9
	v_mul_hi_u32 v15, v13, v9
	v_sub_u32_e32 v16, vcc, 0, v14
	v_cmp_eq_u32_e64 s[0:1], 0, v15
	v_cndmask_b32_e64 v14, v14, v16, s[0:1]
	v_mul_hi_u32 v14, v14, v13
	v_add_u32_e32 v15, vcc, v14, v13
	v_subrev_u32_e32 v13, vcc, v14, v13
	v_cndmask_b32_e64 v13, v13, v15, s[0:1]
	v_mul_hi_u32 v13, v13, v8
	v_mul_lo_u32 v14, v13, v9
	v_add_u32_e32 v15, vcc, 1, v13
	v_add_u32_e32 v16, vcc, -1, v13
	v_sub_u32_e32 v17, vcc, v8, v14
	v_cmp_ge_u32_e32 vcc, v8, v14
	v_cmp_ge_u32_e64 s[0:1], v17, v9
	s_and_b64 s[0:1], s[0:1], vcc
	v_cndmask_b32_e64 v8, v13, v15, s[0:1]
	v_cndmask_b32_e32 v8, v16, v8, vcc
	v_xor_b32_e32 v8, v8, v10
	v_sub_u32_e32 v8, vcc, v8, v10
	v_ashrrev_i32_e32 v9, 31, v8
	v_xor_b32_e32 v10, v11, v9
	v_add_u32_e32 v11, vcc, v9, v8
	v_xor_b32_e32 v9, v11, v9
	v_cvt_f32_u32_e32 v11, v9
	v_rcp_iflag_f32_e32 v11, v11
	v_mul_f32_e32 v11, s5, v11
	v_cvt_u32_f32_e32 v11, v11
	v_mul_lo_u32 v13, v11, v9
	v_mul_hi_u32 v14, v11, v9
	v_sub_u32_e32 v15, vcc, 0, v13
	v_cmp_eq_u32_e64 s[0:1], 0, v14
	v_cndmask_b32_e64 v13, v13, v15, s[0:1]
	v_mul_hi_u32 v13, v13, v11
	v_add_u32_e32 v14, vcc, v13, v11
	v_subrev_u32_e32 v11, vcc, v13, v11
	v_cndmask_b32_e64 v11, v11, v14, s[0:1]
	v_mul_hi_u32 v11, v11, v12
	v_mul_lo_u32 v13, v11, v9
	v_add_u32_e32 v14, vcc, 1, v11
	v_add_u32_e32 v15, vcc, -1, v11
	v_sub_u32_e32 v16, vcc, v12, v13
	v_cmp_ge_u32_e32 vcc, v12, v13
	v_cmp_ge_u32_e64 s[0:1], v16, v9
	s_and_b64 s[0:1], s[0:1], vcc
	v_cndmask_b32_e64 v9, v11, v14, s[0:1]
	v_cndmask_b32_e32 v9, v15, v9, vcc
	v_xor_b32_e32 v9, v9, v10
	v_sub_u32_e32 v9, vcc, v9, v10
	v_mul_lo_u32 v10, v9, v8
	s_add_u32 s2, s2, 4
	s_addc_u32 s3, s3, 0
	flat_store_dword v[6:7], v9
	v_add_u32_e32 v6, vcc, 4, v6
	v_addc_u32_e32 v7, vcc, 0, v7, vcc
	s_cmp_lg_u32 s6, 0
	v_sub_u32_e32 v5, vcc, v5, v10
	s_cbranch_scc1 BB2_4
BB2_5:                                  ; %Flow77
	s_add_i32 s6, s4, -1
	v_lshlrev_b64 v[1:2], 2, v[1:2]
	s_ashr_i32 s7, s6, 31
	s_lshl_b64 s[24:25], s[6:7], 2
	s_add_u32 s26, s8, 4
	v_mov_b32_e32 v5, s11
	v_add_u32_e32 v1, vcc, s10, v1
	v_addc_u32_e32 v2, vcc, v5, v2, vcc
	s_addc_u32 s27, s9, 0
	v_mov_b32_e32 v5, s26
	v_mov_b32_e32 v6, s27
	flat_load_dword v6, v[5:6]
	v_mov_b32_e32 v9, s25
	s_mov_b32 s5, 0x4f800000
	s_waitcnt vmcnt(0) lgkmcnt(0)
	v_ashrrev_i32_e32 v5, 31, v6
	v_add_u32_e32 v6, vcc, v6, v5
	v_add_u32_e32 v8, vcc, s24, v3
	v_addc_u32_e32 v9, vcc, v4, v9, vcc
	flat_load_dword v11, v[8:9]
	v_xor_b32_e32 v6, v6, v5
	v_cvt_f32_u32_e32 v7, v6
	v_rcp_iflag_f32_e32 v7, v7
	v_mul_f32_e32 v7, s5, v7
	v_cvt_u32_f32_e32 v7, v7
	v_mul_lo_u32 v8, v7, v6
	v_mul_hi_u32 v9, v7, v6
	v_sub_u32_e32 v10, vcc, 0, v8
	v_cmp_eq_u32_e64 s[0:1], 0, v9
	v_cndmask_b32_e64 v13, v8, v10, s[0:1]
	v_mul_hi_u32 v13, v13, v7
	s_waitcnt vmcnt(0) lgkmcnt(0)
	v_ashrrev_i32_e32 v12, 31, v11
	v_add_u32_e32 v11, vcc, v11, v12
	v_xor_b32_e32 v14, v11, v12
	v_add_u32_e32 v11, vcc, v13, v7
	v_subrev_u32_e32 v13, vcc, v13, v7
	v_cndmask_b32_e64 v11, v13, v11, s[0:1]
	v_mul_hi_u32 v11, v11, v14
	v_mul_lo_u32 v13, v11, v6
	v_mov_b32_e32 v11, 0
	v_sub_u32_e32 v15, vcc, v14, v13
	v_add_u32_e64 v16, s[0:1], v15, v6
	v_cmp_ge_u32_e64 s[0:1], v14, v13
	v_cmp_ge_u32_e32 vcc, v15, v6
	v_sub_u32_e64 v13, s[2:3], v15, v6
	s_and_b64 vcc, vcc, s[0:1]
	v_cndmask_b32_e32 v13, v15, v13, vcc
	v_cndmask_b32_e64 v13, v16, v13, s[0:1]
	v_xor_b32_e32 v13, v13, v12
	v_sub_u32_e32 v12, vcc, v13, v12
	v_cmp_eq_u32_e32 vcc, 0, v12
	s_and_saveexec_b64 s[10:11], vcc
	s_cbranch_execz BB2_18
; %bb.6:
	v_mov_b32_e32 v12, s9
	v_mov_b32_e32 v11, s8
	flat_load_dword v11, v[11:12]
	s_add_i32 s0, s4, -2
	s_ashr_i32 s1, s0, 31
	s_lshl_b64 s[22:23], s[0:1], 2
	s_waitcnt vmcnt(0) lgkmcnt(0)
	v_ashrrev_i32_e32 v12, 31, v11
	v_add_u32_e32 v11, vcc, v11, v12
	v_xor_b32_e32 v13, v11, v12
	v_cvt_f32_u32_e32 v11, v13
	v_mov_b32_e32 v12, s23
	v_rcp_iflag_f32_e32 v14, v11
	v_add_u32_e32 v11, vcc, s22, v3
	v_addc_u32_e32 v12, vcc, v4, v12, vcc
	flat_load_dword v11, v[11:12]
	v_mul_f32_e32 v14, s5, v14
	v_cvt_u32_f32_e32 v14, v14
	v_mul_lo_u32 v12, v14, v13
	v_mul_hi_u32 v15, v14, v13
	v_sub_u32_e32 v17, vcc, 0, v12
	v_cmp_eq_u32_e64 s[0:1], 0, v15
	v_cndmask_b32_e64 v12, v12, v17, s[0:1]
	v_mul_hi_u32 v12, v12, v14
	v_add_u32_e32 v15, vcc, v12, v14
	v_subrev_u32_e32 v12, vcc, v12, v14
	v_cndmask_b32_e64 v12, v12, v15, s[0:1]
	s_waitcnt vmcnt(0) lgkmcnt(0)
	v_ashrrev_i32_e32 v16, 31, v11
	v_add_u32_e32 v11, vcc, v11, v16
	v_xor_b32_e32 v11, v11, v16
	v_mul_hi_u32 v12, v12, v11
	v_mul_lo_u32 v12, v12, v13
	v_sub_u32_e32 v14, vcc, v11, v12
	v_add_u32_e64 v15, s[0:1], v14, v13
	v_cmp_ge_u32_e64 s[0:1], v11, v12
	v_cmp_ge_u32_e32 vcc, v14, v13
	v_sub_u32_e64 v11, s[2:3], v14, v13
	s_and_b64 vcc, vcc, s[0:1]
	v_cndmask_b32_e32 v11, v14, v11, vcc
	v_cndmask_b32_e64 v11, v15, v11, s[0:1]
	v_xor_b32_e32 v11, v11, v16
	v_sub_u32_e32 v11, vcc, v11, v16
	v_cmp_eq_u32_e32 vcc, 0, v11
	v_mov_b32_e32 v11, 0
	s_and_saveexec_b64 s[18:19], vcc
	s_cbranch_execz BB2_17
; %bb.7:
	s_andn2_b64 vcc, exec, s[20:21]
	s_cbranch_vccnz BB2_11
; %bb.8:
	v_mov_b32_e32 v6, v2
	v_mov_b32_e32 v5, v1
BB2_9:                                  ; %.preheader1
                                        ; =>This Inner Loop Header: Depth=1
	flat_load_dword v7, v[3:4]
	v_add_u32_e32 v3, vcc, 4, v3
	v_addc_u32_e32 v4, vcc, 0, v4, vcc
	s_add_i32 s4, s4, -1
	s_cmp_eq_u32 s4, 0
	s_waitcnt vmcnt(0) lgkmcnt(0)
	flat_store_dword v[5:6], v7
	v_add_u32_e32 v5, vcc, 4, v5
	v_addc_u32_e32 v6, vcc, 0, v6, vcc
	s_cbranch_scc0 BB2_9
; %bb.10:
	v_mov_b32_e32 v3, s26
	v_mov_b32_e32 v4, s27
	flat_load_dword v3, v[3:4]
	s_mov_b64 s[2:3], -1
	s_waitcnt vmcnt(0) lgkmcnt(0)
	v_ashrrev_i32_e32 v5, 31, v3
	v_add_u32_e32 v3, vcc, v3, v5
	v_xor_b32_e32 v6, v3, v5
	v_cvt_f32_u32_e32 v3, v6
	v_rcp_iflag_f32_e32 v3, v3
	v_mul_f32_e32 v3, 0x4f800000, v3
	v_cvt_u32_f32_e32 v7, v3
	v_mul_lo_u32 v8, v7, v6
	v_mul_hi_u32 v9, v7, v6
	v_sub_u32_e32 v10, vcc, 0, v8
	s_branch BB2_12
BB2_11:
	s_mov_b64 s[2:3], 0
BB2_12:
	v_mov_b32_e32 v4, s25
	v_add_u32_e32 v3, vcc, s24, v1
	v_addc_u32_e32 v4, vcc, v2, v4, vcc
	v_cmp_eq_u32_e64 s[0:1], 0, v9
	flat_load_dword v9, v[3:4]
	v_cndmask_b32_e64 v8, v8, v10, s[0:1]
	v_mul_hi_u32 v8, v8, v7
	v_add_u32_e32 v10, vcc, v8, v7
	v_subrev_u32_e32 v7, vcc, v8, v7
	v_cndmask_b32_e64 v7, v7, v10, s[0:1]
	s_waitcnt vmcnt(0) lgkmcnt(0)
	v_ashrrev_i32_e32 v10, 31, v9
	v_add_u32_e32 v8, vcc, v9, v10
	v_xor_b32_e32 v9, v8, v10
	v_mul_hi_u32 v11, v7, v9
	v_xor_b32_e32 v5, v10, v5
	v_mov_b32_e32 v7, s8
	v_mov_b32_e32 v8, s9
	v_mul_lo_u32 v10, v11, v6
	v_add_u32_e32 v12, vcc, 1, v11
	v_add_u32_e32 v13, vcc, -1, v11
	v_sub_u32_e32 v14, vcc, v9, v10
	v_cmp_ge_u32_e32 vcc, v9, v10
	v_cmp_ge_u32_e64 s[0:1], v14, v6
	s_and_b64 s[0:1], s[0:1], vcc
	v_cndmask_b32_e64 v6, v11, v12, s[0:1]
	v_cndmask_b32_e32 v6, v13, v6, vcc
	v_xor_b32_e32 v6, v6, v5
	v_sub_u32_e32 v5, vcc, v6, v5
	flat_store_dword v[3:4], v5
	flat_load_dword v3, v[7:8]
	v_mov_b32_e32 v4, s23
	s_waitcnt vmcnt(0) lgkmcnt(0)
	v_ashrrev_i32_e32 v5, 31, v3
	v_add_u32_e32 v3, vcc, v3, v5
	v_xor_b32_e32 v6, v3, v5
	v_cvt_f32_u32_e32 v3, v6
	v_rcp_iflag_f32_e32 v7, v3
	v_add_u32_e32 v3, vcc, s22, v1
	v_addc_u32_e32 v4, vcc, v2, v4, vcc
	flat_load_dword v8, v[3:4]
	v_mul_f32_e32 v7, 0x4f800000, v7
	v_cvt_u32_f32_e32 v7, v7
	v_mul_lo_u32 v9, v7, v6
	v_mul_hi_u32 v10, v7, v6
	v_sub_u32_e32 v12, vcc, 0, v9
	v_cmp_eq_u32_e64 s[0:1], 0, v10
	v_cndmask_b32_e64 v9, v9, v12, s[0:1]
	v_mul_hi_u32 v9, v9, v7
	v_add_u32_e32 v10, vcc, v9, v7
	v_subrev_u32_e32 v7, vcc, v9, v7
	v_cndmask_b32_e64 v7, v7, v10, s[0:1]
	s_waitcnt vmcnt(0) lgkmcnt(0)
	v_ashrrev_i32_e32 v11, 31, v8
	v_add_u32_e32 v8, vcc, v8, v11
	v_xor_b32_e32 v8, v8, v11
	v_mul_hi_u32 v7, v7, v8
	v_xor_b32_e32 v5, v11, v5
	s_andn2_b64 vcc, exec, s[2:3]
	v_mul_lo_u32 v9, v7, v6
	v_add_u32_e64 v10, s[0:1], 1, v7
	v_add_u32_e64 v11, s[0:1], -1, v7
	v_sub_u32_e64 v12, s[0:1], v8, v9
	v_cmp_ge_u32_e64 s[0:1], v8, v9
	v_cmp_ge_u32_e64 s[2:3], v12, v6
	s_and_b64 s[2:3], s[2:3], s[0:1]
	v_cndmask_b32_e64 v6, v7, v10, s[2:3]
	v_cndmask_b32_e64 v6, v11, v6, s[0:1]
	v_xor_b32_e32 v6, v6, v5
	v_sub_u32_e64 v5, s[0:1], v6, v5
	flat_store_dword v[3:4], v5
	v_mov_b32_e32 v3, 0
	v_mov_b32_e32 v4, 0
	s_cbranch_vccnz BB2_16
; %bb.13:                               ; %.preheader.preheader
	v_mov_b32_e32 v3, 0
	v_mov_b32_e32 v4, 1
BB2_14:                                 ; %.preheader
                                        ; =>This Inner Loop Header: Depth=1
	s_ashr_i32 s7, s6, 31
	s_lshl_b64 s[0:1], s[6:7], 2
	s_add_u32 s2, s16, s0
	v_mov_b32_e32 v6, s1
	v_add_u32_e32 v5, vcc, s0, v1
	v_addc_u32_e32 v6, vcc, v2, v6, vcc
	s_addc_u32 s3, s17, s1
	flat_load_dword v7, v[5:6]
	v_mov_b32_e32 v6, s3
	v_mov_b32_e32 v5, s2
	flat_load_dword v5, v[5:6]
	s_add_i32 s6, s6, -1
	s_cmp_lg_u32 s6, -1
	s_waitcnt vmcnt(1) lgkmcnt(1)
	v_mul_lo_u32 v6, v7, v4
	s_waitcnt vmcnt(0) lgkmcnt(0)
	v_mul_lo_u32 v4, v5, v4
	v_add_u32_e32 v3, vcc, v6, v3
	s_cbranch_scc1 BB2_14
; %bb.15:
	v_ashrrev_i32_e32 v4, 31, v3
BB2_16:
	v_lshlrev_b64 v[1:2], 2, v[3:4]
	v_mov_b32_e32 v3, s13
	v_add_u32_e32 v1, vcc, s12, v1
	v_addc_u32_e32 v2, vcc, v3, v2, vcc
	flat_load_dword v11, v[1:2]
BB2_17:                                 ; %Flow75
	s_or_b64 exec, exec, s[18:19]
BB2_18:
	s_or_b64 exec, exec, s[10:11]
	v_mov_b32_e32 v1, 0
	v_mov_b32_e32 v2, v0
	v_ashrrev_i64 v[0:1], 30, v[1:2]
	v_mov_b32_e32 v2, s15
	v_add_u32_e32 v0, vcc, s14, v0
	v_addc_u32_e32 v1, vcc, v2, v1, vcc
	s_waitcnt vmcnt(0) lgkmcnt(0)
	flat_store_dword v[0:1], v11
	s_endpgm
.Lfunc_end2:
	.size	dilate_tensor, .Lfunc_end2-dilate_tensor
                                        ; -- End function
	.section	.AMDGPU.csdata
; Kernel info:
; codeLenInByte = 1884
; NumSgprs: 30
; NumVgprs: 18
; ScratchSize: 0
; MemoryBound: 0
; FloatMode: 192
; IeeeMode: 1
; LDSByteSize: 0 bytes/workgroup (compile time only)
; SGPRBlocks: 3
; VGPRBlocks: 4
; NumSGPRsForWavesPerEU: 30
; NumVGPRsForWavesPerEU: 18
; Occupancy: 10
; WaveLimiterHint : 1
; COMPUTE_PGM_RSRC2:USER_SGPR: 8
; COMPUTE_PGM_RSRC2:TRAP_HANDLER: 0
; COMPUTE_PGM_RSRC2:TGID_X_EN: 1
; COMPUTE_PGM_RSRC2:TGID_Y_EN: 0
; COMPUTE_PGM_RSRC2:TGID_Z_EN: 0
; COMPUTE_PGM_RSRC2:TIDIG_COMP_CNT: 0
	.text
	.protected	softmax_exp     ; -- Begin function softmax_exp
	.globl	softmax_exp
	.p2align	8
	.type	softmax_exp,@function
	.amdgpu_hsa_kernel softmax_exp
softmax_exp:                            ; @softmax_exp
softmax_exp$local:
	.amd_kernel_code_t
		amd_code_version_major = 1
		amd_code_version_minor = 2
		amd_machine_kind = 1
		amd_machine_version_major = 8
		amd_machine_version_minor = 0
		amd_machine_version_stepping = 3
		kernel_code_entry_byte_offset = 256
		kernel_code_prefetch_byte_size = 0
		granulated_workitem_vgpr_count = 2
		granulated_wavefront_sgpr_count = 1
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
		user_sgpr_count = 8
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
		enable_sgpr_dispatch_ptr = 1
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
		workgroup_group_segment_byte_size = 0
		gds_segment_byte_size = 0
		kernarg_segment_byte_size = 80
		workgroup_fbarrier_count = 0
		wavefront_sgpr_count = 11
		workitem_vgpr_count = 9
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
	s_load_dword s0, s[4:5], 0x4
	s_waitcnt lgkmcnt(0)
	s_and_b32 s0, s0, 0xffff
	s_mul_i32 s8, s8, s0
	s_load_dword s2, s[6:7], 0x10
	s_load_dwordx2 s[0:1], s[6:7], 0x18
	v_add_u32_e32 v0, vcc, s8, v0
	s_waitcnt lgkmcnt(0)
	v_mov_b32_e32 v1, s1
	v_add_u32_e32 v0, vcc, s0, v0
	v_addc_u32_e32 v1, vcc, 0, v1, vcc
	v_cmp_gt_u32_e32 vcc, s2, v0
	s_and_saveexec_b64 s[0:1], vcc
	s_cbranch_execz BB3_2
; %bb.1:
	s_load_dwordx4 s[0:3], s[6:7], 0x0
	v_lshlrev_b64 v[0:1], 2, v[0:1]
	s_mov_b32 s4, 0x42b17218
	v_and_b32_e32 v3, 3, v1
	v_mov_b32_e32 v5, 0x7f800000
	s_waitcnt lgkmcnt(0)
	v_mov_b32_e32 v2, s1
	v_add_u32_e32 v1, vcc, s0, v0
	v_addc_u32_e32 v2, vcc, v2, v3, vcc
	flat_load_dword v2, v[1:2]
	s_mov_b32 s0, 0x39a3b295
	s_mov_b32 s1, 0x3fb8a000
	v_mov_b32_e32 v4, s3
	v_add_u32_e32 v0, vcc, s2, v0
	s_mov_b32 s3, 0xc2aeac50
	s_waitcnt vmcnt(0) lgkmcnt(0)
	v_and_b32_e32 v1, 0xfffff000, v2
	v_sub_f32_e32 v6, v2, v1
	v_mul_f32_e32 v8, s0, v6
	v_mac_f32_e32 v8, s1, v6
	v_mul_f32_e32 v7, s1, v1
	v_mac_f32_e32 v8, s0, v1
	v_exp_f32_e32 v7, v7
	v_exp_f32_e32 v6, v8
	v_addc_u32_e32 v1, vcc, v4, v3, vcc
	v_cmp_ngt_f32_e32 vcc, s3, v2
	v_mul_f32_e32 v3, v7, v6
	v_cndmask_b32_e32 v3, 0, v3, vcc
	v_cmp_nlt_f32_e32 vcc, s4, v2
	v_cndmask_b32_e32 v2, v5, v3, vcc
	flat_store_dword v[0:1], v2
BB3_2:
	s_endpgm
.Lfunc_end3:
	.size	softmax_exp, .Lfunc_end3-softmax_exp
                                        ; -- End function
	.section	.AMDGPU.csdata
; Kernel info:
; codeLenInByte = 240
; NumSgprs: 11
; NumVgprs: 9
; ScratchSize: 0
; MemoryBound: 0
; FloatMode: 192
; IeeeMode: 1
; LDSByteSize: 0 bytes/workgroup (compile time only)
; SGPRBlocks: 1
; VGPRBlocks: 2
; NumSGPRsForWavesPerEU: 11
; NumVGPRsForWavesPerEU: 9
; Occupancy: 10
; WaveLimiterHint : 1
; COMPUTE_PGM_RSRC2:USER_SGPR: 8
; COMPUTE_PGM_RSRC2:TRAP_HANDLER: 0
; COMPUTE_PGM_RSRC2:TGID_X_EN: 1
; COMPUTE_PGM_RSRC2:TGID_Y_EN: 0
; COMPUTE_PGM_RSRC2:TGID_Z_EN: 0
; COMPUTE_PGM_RSRC2:TIDIG_COMP_CNT: 0
	.text
	.protected	softmax_div     ; -- Begin function softmax_div
	.globl	softmax_div
	.p2align	8
	.type	softmax_div,@function
	.amdgpu_hsa_kernel softmax_div
softmax_div:                            ; @softmax_div
softmax_div$local:
	.amd_kernel_code_t
		amd_code_version_major = 1
		amd_code_version_minor = 2
		amd_machine_kind = 1
		amd_machine_version_major = 8
		amd_machine_version_minor = 0
		amd_machine_version_stepping = 3
		kernel_code_entry_byte_offset = 256
		kernel_code_prefetch_byte_size = 0
		granulated_workitem_vgpr_count = 2
		granulated_wavefront_sgpr_count = 1
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
		user_sgpr_count = 8
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
		enable_sgpr_dispatch_ptr = 1
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
		workgroup_group_segment_byte_size = 0
		gds_segment_byte_size = 0
		kernarg_segment_byte_size = 88
		workgroup_fbarrier_count = 0
		wavefront_sgpr_count = 14
		workitem_vgpr_count = 10
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
	s_load_dwordx2 s[2:3], s[6:7], 0x18
	s_load_dwordx2 s[0:1], s[6:7], 0x20
	s_waitcnt lgkmcnt(0)
	s_load_dword s1, s[4:5], 0x4
	s_waitcnt lgkmcnt(0)
	s_and_b32 s1, s1, 0xffff
	s_mul_i32 s8, s8, s1
	v_add_u32_e32 v0, vcc, s8, v0
	v_add_u32_e32 v0, vcc, s0, v0
	v_cmp_ge_i32_e32 vcc, s2, v0
	s_and_saveexec_b64 s[0:1], vcc
	s_cbranch_execz BB4_2
; %bb.1:
	s_ashr_i32 s4, s3, 31
	s_add_i32 s0, s3, s4
	s_xor_b32 s5, s0, s4
	v_cvt_f32_u32_e32 v1, s5
	s_mov_b32 s8, 0x4f800000
	s_ashr_i32 s3, s2, 31
	s_add_i32 s2, s2, s3
	v_rcp_iflag_f32_e32 v1, v1
	s_xor_b32 s2, s2, s3
	s_xor_b32 s3, s3, s4
	v_ashrrev_i32_e32 v6, 31, v0
	v_mul_f32_e32 v1, s8, v1
	v_cvt_u32_f32_e32 v1, v1
	v_mul_lo_u32 v2, v1, s5
	v_mul_hi_u32 v3, v1, s5
	v_sub_u32_e32 v4, vcc, 0, v2
	v_cmp_eq_u32_e64 s[0:1], 0, v3
	v_cndmask_b32_e64 v2, v2, v4, s[0:1]
	v_mul_hi_u32 v2, v2, v1
	v_add_u32_e32 v3, vcc, v2, v1
	v_subrev_u32_e32 v1, vcc, v2, v1
	v_cndmask_b32_e64 v1, v1, v3, s[0:1]
	v_mul_hi_u32 v1, v1, s2
	v_mul_lo_u32 v2, v1, s5
	v_add_u32_e32 v3, vcc, 1, v1
	v_add_u32_e32 v4, vcc, -1, v1
	v_sub_u32_e32 v5, vcc, s2, v2
	v_cmp_ge_u32_e32 vcc, s2, v2
	v_cmp_le_u32_e64 s[0:1], s5, v5
	s_and_b64 s[0:1], s[0:1], vcc
	v_cndmask_b32_e64 v1, v1, v3, s[0:1]
	v_cndmask_b32_e32 v1, v4, v1, vcc
	v_xor_b32_e32 v1, s3, v1
	v_subrev_u32_e32 v1, vcc, s3, v1
	v_ashrrev_i32_e32 v2, 31, v1
	v_add_u32_e32 v1, vcc, v2, v1
	v_xor_b32_e32 v3, v1, v2
	v_cvt_f32_u32_e32 v1, v3
	v_xor_b32_e32 v2, v6, v2
	v_rcp_iflag_f32_e32 v1, v1
	v_mul_f32_e32 v1, s8, v1
	v_cvt_u32_f32_e32 v1, v1
	s_load_dwordx4 s[8:11], s[6:7], 0x0
	s_load_dwordx4 s[4:7], s[6:7], 0x10
	v_mul_lo_u32 v4, v1, v3
	v_mul_hi_u32 v5, v1, v3
	v_sub_u32_e32 v7, vcc, 0, v4
	v_cmp_eq_u32_e64 s[0:1], 0, v5
	v_cndmask_b32_e64 v4, v4, v7, s[0:1]
	v_mul_hi_u32 v4, v4, v1
	v_add_u32_e32 v5, vcc, v0, v6
	v_xor_b32_e32 v5, v5, v6
	v_add_u32_e32 v7, vcc, v4, v1
	v_subrev_u32_e32 v1, vcc, v4, v1
	v_cndmask_b32_e64 v1, v1, v7, s[0:1]
	v_mul_hi_u32 v4, v1, v5
	v_mov_b32_e32 v1, 0
	v_mul_lo_u32 v6, v4, v3
	v_add_u32_e32 v7, vcc, 1, v4
	v_add_u32_e32 v8, vcc, -1, v4
	v_sub_u32_e32 v9, vcc, v5, v6
	v_cmp_ge_u32_e32 vcc, v5, v6
	v_cmp_ge_u32_e64 s[0:1], v9, v3
	s_and_b64 s[0:1], s[0:1], vcc
	v_cndmask_b32_e64 v3, v4, v7, s[0:1]
	v_cndmask_b32_e32 v3, v8, v3, vcc
	v_xor_b32_e32 v3, v3, v2
	v_sub_u32_e32 v3, vcc, v3, v2
	v_mov_b32_e32 v2, v0
	v_ashrrev_i64 v[0:1], 30, v[1:2]
	s_waitcnt lgkmcnt(0)
	v_mov_b32_e32 v2, s9
	v_add_u32_e32 v5, vcc, s8, v0
	v_ashrrev_i32_e32 v4, 31, v3
	v_addc_u32_e32 v6, vcc, v2, v1, vcc
	v_lshlrev_b64 v[2:3], 2, v[3:4]
	v_mov_b32_e32 v4, s5
	v_add_u32_e32 v2, vcc, s4, v2
	v_addc_u32_e32 v3, vcc, v4, v3, vcc
	flat_load_dword v2, v[2:3]
	s_mov_b32 s0, 0x6f800000
	v_mov_b32_e32 v4, 0x2f800000
	v_mov_b32_e32 v7, s11
	s_waitcnt vmcnt(0) lgkmcnt(0)
	v_cmp_gt_f32_e64 vcc, |v2|, s0
	v_cndmask_b32_e32 v3, 1.0, v4, vcc
	flat_load_dword v4, v[5:6]
	v_mul_f32_e32 v2, v2, v3
	v_rcp_f32_e32 v2, v2
	v_add_u32_e32 v0, vcc, s10, v0
	v_addc_u32_e32 v1, vcc, v7, v1, vcc
	s_waitcnt vmcnt(0) lgkmcnt(0)
	v_mul_f32_e32 v2, v4, v2
	v_mul_f32_e32 v2, v3, v2
	flat_store_dword v[0:1], v2
BB4_2:
	s_endpgm
.Lfunc_end4:
	.size	softmax_div, .Lfunc_end4-softmax_div
                                        ; -- End function
	.section	.AMDGPU.csdata
; Kernel info:
; codeLenInByte = 572
; NumSgprs: 14
; NumVgprs: 10
; ScratchSize: 0
; MemoryBound: 0
; FloatMode: 192
; IeeeMode: 1
; LDSByteSize: 0 bytes/workgroup (compile time only)
; SGPRBlocks: 1
; VGPRBlocks: 2
; NumSGPRsForWavesPerEU: 14
; NumVGPRsForWavesPerEU: 10
; Occupancy: 10
; WaveLimiterHint : 1
; COMPUTE_PGM_RSRC2:USER_SGPR: 8
; COMPUTE_PGM_RSRC2:TRAP_HANDLER: 0
; COMPUTE_PGM_RSRC2:TGID_X_EN: 1
; COMPUTE_PGM_RSRC2:TGID_Y_EN: 0
; COMPUTE_PGM_RSRC2:TGID_Z_EN: 0
; COMPUTE_PGM_RSRC2:TIDIG_COMP_CNT: 0
	.text
	.protected	sum_one_axis    ; -- Begin function sum_one_axis
	.globl	sum_one_axis
	.p2align	8
	.type	sum_one_axis,@function
	.amdgpu_hsa_kernel sum_one_axis
sum_one_axis:                           ; @sum_one_axis
sum_one_axis$local:
	.amd_kernel_code_t
		amd_code_version_major = 1
		amd_code_version_minor = 2
		amd_machine_kind = 1
		amd_machine_version_major = 8
		amd_machine_version_minor = 0
		amd_machine_version_stepping = 3
		kernel_code_entry_byte_offset = 256
		kernel_code_prefetch_byte_size = 0
		granulated_workitem_vgpr_count = 4
		granulated_wavefront_sgpr_count = 3
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
		user_sgpr_count = 8
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
		enable_sgpr_dispatch_ptr = 1
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
		workgroup_group_segment_byte_size = 0
		gds_segment_byte_size = 0
		kernarg_segment_byte_size = 112
		workgroup_fbarrier_count = 0
		wavefront_sgpr_count = 26
		workitem_vgpr_count = 17
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
	s_load_dwordx8 s[12:19], s[6:7], 0x0
	s_load_dwordx2 s[2:3], s[6:7], 0x20
	s_load_dwordx4 s[20:23], s[6:7], 0x28
	s_load_dwordx2 s[0:1], s[6:7], 0x38
	s_waitcnt lgkmcnt(0)
	s_load_dword s1, s[4:5], 0x4
	s_add_i32 s4, s2, -1
	v_mov_b32_e32 v5, s23
	s_waitcnt lgkmcnt(0)
	s_and_b32 s1, s1, 0xffff
	s_mul_i32 s8, s8, s1
	v_add_u32_e32 v0, vcc, s8, v0
	v_add_u32_e32 v0, vcc, s0, v0
	v_mul_lo_u32 v3, v0, s2
	s_cmp_lt_i32 s2, 2
	v_sub_u32_e32 v1, vcc, v3, v0
	v_ashrrev_i32_e32 v2, 31, v1
	v_lshlrev_b64 v[1:2], 2, v[1:2]
	v_add_u32_e32 v1, vcc, s22, v1
	v_addc_u32_e32 v2, vcc, v5, v2, vcc
	s_cbranch_scc1 BB5_5
; %bb.1:                                ; %.preheader9.preheader
	v_mov_b32_e32 v4, v0
	s_mov_b64 s[6:7], s[18:19]
	s_mov_b32 s0, 1
	s_mov_b64 s[8:9], s[18:19]
	s_mov_b32 s1, s4
BB5_2:                                  ; %.preheader9
                                        ; =>This Inner Loop Header: Depth=1
	s_load_dword s5, s[8:9], 0x0
	s_add_i32 s1, s1, -1
	s_add_u32 s8, s8, 4
	s_addc_u32 s9, s9, 0
	s_cmp_lg_u32 s1, 0
	s_waitcnt lgkmcnt(0)
	s_mul_i32 s0, s5, s0
	s_cbranch_scc1 BB5_2
; %bb.3:                                ; %.preheader7.preheader
	v_mov_b32_e32 v6, v2
	v_mov_b32_e32 v7, s0
	s_mov_b32 s5, 0x4f800000
	v_mov_b32_e32 v5, v1
	s_mov_b32 s8, s4
BB5_4:                                  ; %.preheader7
                                        ; =>This Inner Loop Header: Depth=1
	v_mov_b32_e32 v9, s7
	v_mov_b32_e32 v8, s6
	flat_load_dword v8, v[8:9]
	v_ashrrev_i32_e32 v9, 31, v7
	v_add_u32_e32 v7, vcc, v9, v7
	v_xor_b32_e32 v7, v7, v9
	v_ashrrev_i32_e32 v10, 31, v4
	v_add_u32_e32 v11, vcc, v4, v10
	v_xor_b32_e32 v11, v11, v10
	s_add_i32 s8, s8, -1
	s_waitcnt vmcnt(0) lgkmcnt(0)
	v_ashrrev_i32_e32 v12, 31, v8
	v_add_u32_e32 v8, vcc, v8, v12
	v_xor_b32_e32 v8, v8, v12
	v_xor_b32_e32 v9, v9, v12
	v_cvt_f32_u32_e32 v12, v8
	v_rcp_iflag_f32_e32 v12, v12
	v_mul_f32_e32 v12, s5, v12
	v_cvt_u32_f32_e32 v12, v12
	v_mul_lo_u32 v13, v12, v8
	v_mul_hi_u32 v14, v12, v8
	v_sub_u32_e32 v15, vcc, 0, v13
	v_cmp_eq_u32_e64 s[0:1], 0, v14
	v_cndmask_b32_e64 v13, v13, v15, s[0:1]
	v_mul_hi_u32 v13, v13, v12
	v_add_u32_e32 v14, vcc, v13, v12
	v_subrev_u32_e32 v12, vcc, v13, v12
	v_cndmask_b32_e64 v12, v12, v14, s[0:1]
	v_mul_hi_u32 v12, v12, v7
	v_mul_lo_u32 v13, v12, v8
	v_add_u32_e32 v14, vcc, 1, v12
	v_add_u32_e32 v15, vcc, -1, v12
	v_sub_u32_e32 v16, vcc, v7, v13
	v_cmp_ge_u32_e32 vcc, v7, v13
	v_cmp_ge_u32_e64 s[0:1], v16, v8
	s_and_b64 s[0:1], s[0:1], vcc
	v_cndmask_b32_e64 v7, v12, v14, s[0:1]
	v_cndmask_b32_e32 v7, v15, v7, vcc
	v_xor_b32_e32 v7, v7, v9
	v_sub_u32_e32 v7, vcc, v7, v9
	v_ashrrev_i32_e32 v8, 31, v7
	v_xor_b32_e32 v9, v10, v8
	v_add_u32_e32 v10, vcc, v8, v7
	v_xor_b32_e32 v8, v10, v8
	v_cvt_f32_u32_e32 v10, v8
	v_rcp_iflag_f32_e32 v10, v10
	v_mul_f32_e32 v10, s5, v10
	v_cvt_u32_f32_e32 v10, v10
	v_mul_lo_u32 v12, v10, v8
	v_mul_hi_u32 v13, v10, v8
	v_sub_u32_e32 v14, vcc, 0, v12
	v_cmp_eq_u32_e64 s[0:1], 0, v13
	v_cndmask_b32_e64 v12, v12, v14, s[0:1]
	v_mul_hi_u32 v12, v12, v10
	v_add_u32_e32 v13, vcc, v12, v10
	v_subrev_u32_e32 v10, vcc, v12, v10
	v_cndmask_b32_e64 v10, v10, v13, s[0:1]
	v_mul_hi_u32 v10, v10, v11
	v_mul_lo_u32 v12, v10, v8
	v_add_u32_e32 v13, vcc, 1, v10
	v_add_u32_e32 v14, vcc, -1, v10
	v_sub_u32_e32 v15, vcc, v11, v12
	v_cmp_ge_u32_e32 vcc, v11, v12
	v_cmp_ge_u32_e64 s[0:1], v15, v8
	s_and_b64 s[0:1], s[0:1], vcc
	v_cndmask_b32_e64 v8, v10, v13, s[0:1]
	v_cndmask_b32_e32 v8, v14, v8, vcc
	v_xor_b32_e32 v8, v8, v9
	v_sub_u32_e32 v8, vcc, v8, v9
	v_mul_lo_u32 v9, v8, v7
	s_add_u32 s6, s6, 4
	s_addc_u32 s7, s7, 0
	flat_store_dword v[5:6], v8
	v_add_u32_e32 v5, vcc, 4, v5
	v_addc_u32_e32 v6, vcc, 0, v6, vcc
	s_cmp_eq_u32 s8, 0
	v_sub_u32_e32 v4, vcc, v4, v9
	s_cbranch_scc0 BB5_4
BB5_5:                                  ; %Flow82
	v_ashrrev_i32_e32 v4, 31, v3
	s_ashr_i32 s1, s3, 31
	s_mov_b32 s0, s3
	s_lshl_b64 s[0:1], s[0:1], 2
	v_lshlrev_b64 v[3:4], 2, v[3:4]
	s_add_u32 s0, s16, s0
	v_mov_b32_e32 v5, s21
	v_add_u32_e32 v3, vcc, s20, v3
	s_addc_u32 s1, s17, s1
	v_addc_u32_e32 v4, vcc, v5, v4, vcc
	v_mov_b32_e32 v6, s1
	v_mov_b32_e32 v5, s0
	flat_load_dword v6, v[5:6]
	v_mov_b32_e32 v5, 0
	v_mov_b32_e32 v7, 0
	s_waitcnt vmcnt(0) lgkmcnt(0)
	v_cmp_gt_i32_e32 vcc, 1, v6
	s_and_b64 vcc, exec, vcc
	s_cbranch_vccnz BB5_25
; %bb.6:
	s_cmp_gt_i32 s2, 0
	s_cbranch_scc1 BB5_10
; %bb.7:
	v_mov_b32_e32 v7, s12
	v_mov_b32_e32 v8, s13
	flat_load_dword v8, v[7:8]
	v_mov_b32_e32 v7, 0
BB5_8:                                  ; =>This Inner Loop Header: Depth=1
	v_add_u32_e32 v6, vcc, -1, v6
	v_cmp_eq_u32_e32 vcc, 0, v6
	s_and_b64 vcc, exec, vcc
	s_waitcnt vmcnt(0) lgkmcnt(0)
	v_add_f32_e32 v7, v8, v7
	s_cbranch_vccz BB5_8
; %bb.9:                                ; %Flow
	s_mov_b64 s[6:7], 0
	s_and_b64 vcc, exec, s[6:7]
	s_cbranch_vccnz BB5_11
	s_branch BB5_25
BB5_10:
                                        ; implicit-def: $vgpr7
	s_cbranch_execz BB5_25
BB5_11:
	s_mov_b32 s7, 0
	s_mov_b32 s6, s3
	s_lshl_b64 s[8:9], s[6:7], 2
	v_mov_b32_e32 v6, s9
	v_add_u32_e32 v8, vcc, s8, v3
	v_addc_u32_e32 v9, vcc, v4, v6, vcc
	v_mov_b32_e32 v7, 0
BB5_12:                                 ; =>This Loop Header: Depth=1
                                        ;     Child Loop BB5_14 Depth 2
                                        ;     Child Loop BB5_23 Depth 2
	v_mov_b32_e32 v11, v2
	v_mov_b32_e32 v13, v4
	v_mov_b32_e32 v10, v1
	v_mov_b32_e32 v12, v3
	s_mov_b32 s5, 0
	s_mov_b32 s6, 0
	s_branch BB5_14
BB5_13:                                 ;   in Loop: Header=BB5_14 Depth=2
	v_add_u32_e32 v12, vcc, 4, v12
	v_addc_u32_e32 v13, vcc, 0, v13, vcc
	s_add_i32 s5, s5, 1
	v_add_u32_e32 v10, vcc, 4, v10
	s_cmp_lg_u32 s2, s5
	v_addc_u32_e32 v11, vcc, 0, v11, vcc
	s_cbranch_scc0 BB5_22
BB5_14:                                 ;   Parent Loop BB5_12 Depth=1
                                        ; =>  This Inner Loop Header: Depth=2
	s_cmp_eq_u32 s3, s5
	s_mov_b64 s[8:9], -1
	s_cbranch_scc1 BB5_20
; %bb.15:                               ;   in Loop: Header=BB5_14 Depth=2
	s_cmp_eq_u32 s6, 0
	s_mov_b64 s[8:9], -1
	s_cbranch_scc1 BB5_17
; %bb.16:                               ;   in Loop: Header=BB5_14 Depth=2
	v_add_u32_e32 v14, vcc, -4, v10
	v_addc_u32_e32 v15, vcc, -1, v11, vcc
	flat_load_dword v6, v[14:15]
	s_mov_b64 s[8:9], 0
	s_waitcnt vmcnt(0) lgkmcnt(0)
	flat_store_dword v[12:13], v6
BB5_17:                                 ; %Flow75
                                        ;   in Loop: Header=BB5_14 Depth=2
	s_andn2_b64 vcc, exec, s[8:9]
	s_cbranch_vccnz BB5_19
; %bb.18:                               ;   in Loop: Header=BB5_14 Depth=2
	flat_load_dword v6, v[10:11]
	s_mov_b32 s6, 0
	s_waitcnt vmcnt(0) lgkmcnt(0)
	flat_store_dword v[12:13], v6
BB5_19:                                 ; %Flow76
                                        ;   in Loop: Header=BB5_14 Depth=2
	s_mov_b64 s[8:9], 0
BB5_20:                                 ; %Flow77
                                        ;   in Loop: Header=BB5_14 Depth=2
	s_andn2_b64 vcc, exec, s[8:9]
	s_cbranch_vccnz BB5_13
; %bb.21:                               ;   in Loop: Header=BB5_14 Depth=2
	v_mov_b32_e32 v6, s7
	s_mov_b32 s6, 1
	flat_store_dword v[8:9], v6
	s_branch BB5_13
BB5_22:                                 ; %.preheader.preheader
                                        ;   in Loop: Header=BB5_12 Depth=1
	v_mov_b32_e32 v10, 0
	v_mov_b32_e32 v6, 1
	s_mov_b32 s8, s4
BB5_23:                                 ; %.preheader
                                        ;   Parent Loop BB5_12 Depth=1
                                        ; =>  This Inner Loop Header: Depth=2
	s_ashr_i32 s9, s8, 31
	s_lshl_b64 s[10:11], s[8:9], 2
	v_mov_b32_e32 v12, s11
	v_add_u32_e32 v11, vcc, s10, v3
	s_add_u32 s18, s16, s10
	v_addc_u32_e32 v12, vcc, v4, v12, vcc
	s_addc_u32 s19, s17, s11
	flat_load_dword v13, v[11:12]
	v_mov_b32_e32 v11, s18
	v_mov_b32_e32 v12, s19
	flat_load_dword v11, v[11:12]
	s_add_i32 s8, s8, -1
	s_cmp_eq_u32 s8, -1
	s_waitcnt vmcnt(1) lgkmcnt(1)
	v_mul_lo_u32 v12, v13, v6
	s_waitcnt vmcnt(0) lgkmcnt(0)
	v_mul_lo_u32 v6, v11, v6
	v_add_u32_e32 v10, vcc, v12, v10
	s_cbranch_scc0 BB5_23
; %bb.24:                               ;   in Loop: Header=BB5_12 Depth=1
	v_ashrrev_i32_e32 v11, 31, v10
	v_lshlrev_b64 v[10:11], 2, v[10:11]
	v_mov_b32_e32 v6, s13
	v_add_u32_e32 v10, vcc, s12, v10
	v_addc_u32_e32 v11, vcc, v6, v11, vcc
	flat_load_dword v6, v[10:11]
	v_mov_b32_e32 v11, s1
	v_mov_b32_e32 v10, s0
	flat_load_dword v10, v[10:11]
	s_add_i32 s7, s7, 1
	s_waitcnt vmcnt(1) lgkmcnt(1)
	v_add_f32_e32 v7, v7, v6
	s_waitcnt vmcnt(0) lgkmcnt(0)
	v_cmp_ge_i32_e32 vcc, s7, v10
	s_and_b64 vcc, exec, vcc
	s_cbranch_vccz BB5_12
BB5_25:                                 ; %.loopexit
	v_mov_b32_e32 v6, v0
	v_ashrrev_i64 v[0:1], 30, v[5:6]
	v_mov_b32_e32 v2, s15
	v_add_u32_e32 v0, vcc, s14, v0
	v_addc_u32_e32 v1, vcc, v2, v1, vcc
	flat_store_dword v[0:1], v7
	s_endpgm
.Lfunc_end5:
	.size	sum_one_axis, .Lfunc_end5-sum_one_axis
                                        ; -- End function
	.section	.AMDGPU.csdata
; Kernel info:
; codeLenInByte = 1176
; NumSgprs: 26
; NumVgprs: 17
; ScratchSize: 0
; MemoryBound: 0
; FloatMode: 192
; IeeeMode: 1
; LDSByteSize: 0 bytes/workgroup (compile time only)
; SGPRBlocks: 3
; VGPRBlocks: 4
; NumSGPRsForWavesPerEU: 26
; NumVGPRsForWavesPerEU: 17
; Occupancy: 10
; WaveLimiterHint : 1
; COMPUTE_PGM_RSRC2:USER_SGPR: 8
; COMPUTE_PGM_RSRC2:TRAP_HANDLER: 0
; COMPUTE_PGM_RSRC2:TGID_X_EN: 1
; COMPUTE_PGM_RSRC2:TGID_Y_EN: 0
; COMPUTE_PGM_RSRC2:TGID_Z_EN: 0
; COMPUTE_PGM_RSRC2:TIDIG_COMP_CNT: 0
	.text
	.protected	scaleAdd        ; -- Begin function scaleAdd
	.globl	scaleAdd
	.p2align	8
	.type	scaleAdd,@function
	.amdgpu_hsa_kernel scaleAdd
scaleAdd:                               ; @scaleAdd
scaleAdd$local:
	.amd_kernel_code_t
		amd_code_version_major = 1
		amd_code_version_minor = 2
		amd_machine_kind = 1
		amd_machine_version_major = 8
		amd_machine_version_minor = 0
		amd_machine_version_stepping = 3
		kernel_code_entry_byte_offset = 256
		kernel_code_prefetch_byte_size = 0
		granulated_workitem_vgpr_count = 1
		granulated_wavefront_sgpr_count = 1
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
		user_sgpr_count = 8
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
		enable_sgpr_dispatch_ptr = 1
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
		workgroup_group_segment_byte_size = 0
		gds_segment_byte_size = 0
		kernarg_segment_byte_size = 96
		workgroup_fbarrier_count = 0
		wavefront_sgpr_count = 14
		workitem_vgpr_count = 6
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
	s_load_dwordx4 s[0:3], s[6:7], 0x18
	s_load_dwordx2 s[10:11], s[6:7], 0x28
	s_waitcnt lgkmcnt(0)
	s_load_dword s3, s[4:5], 0x4
	s_waitcnt lgkmcnt(0)
	s_and_b32 s3, s3, 0xffff
	s_mul_i32 s8, s8, s3
	v_add_u32_e32 v0, vcc, s8, v0
	v_add_u32_e32 v0, vcc, s10, v0
	v_cmp_ge_i32_e32 vcc, s2, v0
	s_and_saveexec_b64 s[2:3], vcc
	s_cbranch_execz BB6_2
; %bb.1:
	s_load_dwordx4 s[8:11], s[6:7], 0x0
	s_load_dwordx4 s[4:7], s[6:7], 0x10
	v_mov_b32_e32 v1, 0
	v_mov_b32_e32 v2, v0
	v_ashrrev_i64 v[0:1], 30, v[1:2]
	s_waitcnt lgkmcnt(0)
	v_mov_b32_e32 v3, s11
	v_add_u32_e32 v2, vcc, s10, v0
	v_addc_u32_e32 v3, vcc, v3, v1, vcc
	v_mov_b32_e32 v5, s5
	v_add_u32_e32 v4, vcc, s4, v0
	v_addc_u32_e32 v5, vcc, v5, v1, vcc
	flat_load_dword v4, v[4:5]
	flat_load_dword v2, v[2:3]
	v_mov_b32_e32 v5, s9
	v_add_u32_e32 v0, vcc, s8, v0
	v_addc_u32_e32 v1, vcc, v5, v1, vcc
	s_waitcnt vmcnt(1) lgkmcnt(1)
	v_mul_f32_e32 v3, s1, v4
	s_waitcnt vmcnt(0) lgkmcnt(0)
	v_mac_f32_e32 v3, s0, v2
	flat_store_dword v[0:1], v3
BB6_2:
	s_endpgm
.Lfunc_end6:
	.size	scaleAdd, .Lfunc_end6-scaleAdd
                                        ; -- End function
	.section	.AMDGPU.csdata
; Kernel info:
; codeLenInByte = 180
; NumSgprs: 14
; NumVgprs: 6
; ScratchSize: 0
; MemoryBound: 0
; FloatMode: 192
; IeeeMode: 1
; LDSByteSize: 0 bytes/workgroup (compile time only)
; SGPRBlocks: 1
; VGPRBlocks: 1
; NumSGPRsForWavesPerEU: 14
; NumVGPRsForWavesPerEU: 6
; Occupancy: 10
; WaveLimiterHint : 1
; COMPUTE_PGM_RSRC2:USER_SGPR: 8
; COMPUTE_PGM_RSRC2:TRAP_HANDLER: 0
; COMPUTE_PGM_RSRC2:TGID_X_EN: 1
; COMPUTE_PGM_RSRC2:TGID_Y_EN: 0
; COMPUTE_PGM_RSRC2:TGID_Z_EN: 0
; COMPUTE_PGM_RSRC2:TIDIG_COMP_CNT: 0
	.text
	.protected	mul             ; -- Begin function mul
	.globl	mul
	.p2align	8
	.type	mul,@function
	.amdgpu_hsa_kernel mul
mul:                                    ; @mul
mul$local:
	.amd_kernel_code_t
		amd_code_version_major = 1
		amd_code_version_minor = 2
		amd_machine_kind = 1
		amd_machine_version_major = 8
		amd_machine_version_minor = 0
		amd_machine_version_stepping = 3
		kernel_code_entry_byte_offset = 256
		kernel_code_prefetch_byte_size = 0
		granulated_workitem_vgpr_count = 1
		granulated_wavefront_sgpr_count = 1
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
		user_sgpr_count = 8
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
		enable_sgpr_dispatch_ptr = 1
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
		workgroup_group_segment_byte_size = 0
		gds_segment_byte_size = 0
		kernarg_segment_byte_size = 88
		workgroup_fbarrier_count = 0
		wavefront_sgpr_count = 11
		workitem_vgpr_count = 6
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
	s_load_dword s0, s[4:5], 0x4
	s_waitcnt lgkmcnt(0)
	s_and_b32 s0, s0, 0xffff
	s_mul_i32 s8, s8, s0
	s_load_dword s2, s[6:7], 0x18
	s_load_dwordx2 s[0:1], s[6:7], 0x20
	v_add_u32_e32 v0, vcc, s8, v0
	s_waitcnt lgkmcnt(0)
	v_add_u32_e32 v0, vcc, s0, v0
	v_cmp_ge_i32_e32 vcc, s2, v0
	s_and_saveexec_b64 s[0:1], vcc
	s_cbranch_execz BB7_2
; %bb.1:
	s_load_dwordx4 s[0:3], s[6:7], 0x0
	s_load_dwordx4 s[4:7], s[6:7], 0x10
	v_mov_b32_e32 v1, 0
	v_mov_b32_e32 v2, v0
	v_ashrrev_i64 v[0:1], 30, v[1:2]
	s_waitcnt lgkmcnt(0)
	v_mov_b32_e32 v3, s3
	v_add_u32_e32 v2, vcc, s2, v0
	v_addc_u32_e32 v3, vcc, v3, v1, vcc
	flat_load_dword v4, v[2:3]
	v_mov_b32_e32 v3, s5
	v_add_u32_e32 v2, vcc, s4, v0
	v_addc_u32_e32 v3, vcc, v3, v1, vcc
	flat_load_dword v2, v[2:3]
	v_mov_b32_e32 v5, s1
	v_add_u32_e32 v0, vcc, s0, v0
	v_addc_u32_e32 v1, vcc, v5, v1, vcc
	s_waitcnt vmcnt(0) lgkmcnt(0)
	v_mul_f32_e32 v2, v4, v2
	flat_store_dword v[0:1], v2
BB7_2:
	s_endpgm
.Lfunc_end7:
	.size	mul, .Lfunc_end7-mul
                                        ; -- End function
	.section	.AMDGPU.csdata
; Kernel info:
; codeLenInByte = 172
; NumSgprs: 11
; NumVgprs: 6
; ScratchSize: 0
; MemoryBound: 0
; FloatMode: 192
; IeeeMode: 1
; LDSByteSize: 0 bytes/workgroup (compile time only)
; SGPRBlocks: 1
; VGPRBlocks: 1
; NumSGPRsForWavesPerEU: 11
; NumVGPRsForWavesPerEU: 6
; Occupancy: 10
; WaveLimiterHint : 1
; COMPUTE_PGM_RSRC2:USER_SGPR: 8
; COMPUTE_PGM_RSRC2:TRAP_HANDLER: 0
; COMPUTE_PGM_RSRC2:TGID_X_EN: 1
; COMPUTE_PGM_RSRC2:TGID_Y_EN: 0
; COMPUTE_PGM_RSRC2:TGID_Z_EN: 0
; COMPUTE_PGM_RSRC2:TIDIG_COMP_CNT: 0
	.text
	.protected	rmsProp         ; -- Begin function rmsProp
	.globl	rmsProp
	.p2align	8
	.type	rmsProp,@function
	.amdgpu_hsa_kernel rmsProp
rmsProp:                                ; @rmsProp
rmsProp$local:
	.amd_kernel_code_t
		amd_code_version_major = 1
		amd_code_version_minor = 2
		amd_machine_kind = 1
		amd_machine_version_major = 8
		amd_machine_version_minor = 0
		amd_machine_version_stepping = 3
		kernel_code_entry_byte_offset = 256
		kernel_code_prefetch_byte_size = 0
		granulated_workitem_vgpr_count = 2
		granulated_wavefront_sgpr_count = 1
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
		user_sgpr_count = 8
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
		enable_sgpr_dispatch_ptr = 1
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
		workgroup_group_segment_byte_size = 0
		gds_segment_byte_size = 0
		kernarg_segment_byte_size = 96
		workgroup_fbarrier_count = 0
		wavefront_sgpr_count = 14
		workitem_vgpr_count = 11
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
	s_load_dwordx4 s[0:3], s[6:7], 0x18
	s_load_dwordx2 s[10:11], s[6:7], 0x28
	s_waitcnt lgkmcnt(0)
	s_load_dword s3, s[4:5], 0x4
	s_waitcnt lgkmcnt(0)
	s_and_b32 s3, s3, 0xffff
	s_mul_i32 s8, s8, s3
	v_add_u32_e32 v0, vcc, s8, v0
	v_add_u32_e32 v0, vcc, s10, v0
	v_cmp_ge_i32_e32 vcc, s2, v0
	s_and_saveexec_b64 s[2:3], vcc
	s_cbranch_execz BB8_2
; %bb.1:
	s_load_dwordx4 s[8:11], s[6:7], 0x0
	s_load_dwordx4 s[4:7], s[6:7], 0x10
	v_mov_b32_e32 v1, 0
	v_mov_b32_e32 v2, v0
	v_ashrrev_i64 v[0:1], 30, v[1:2]
	s_waitcnt lgkmcnt(0)
	v_mov_b32_e32 v5, s11
	v_mov_b32_e32 v3, s5
	v_add_u32_e32 v2, vcc, s4, v0
	v_addc_u32_e32 v3, vcc, v3, v1, vcc
	v_add_u32_e32 v4, vcc, s10, v0
	v_addc_u32_e32 v5, vcc, v5, v1, vcc
	flat_load_dword v8, v[4:5]
	flat_load_dword v6, v[2:3]
	v_sub_f32_e64 v7, 1.0, s0
	s_mov_b32 s3, 0x3eb0c6f7
	s_mov_b32 s2, 0xa0b5ed8d
	v_mov_b32_e32 v10, s9
	v_add_u32_e32 v0, vcc, s8, v0
	v_addc_u32_e32 v1, vcc, v10, v1, vcc
	v_mov_b32_e32 v9, 0x2f800000
	s_waitcnt vmcnt(1) lgkmcnt(1)
	v_mul_f32_e32 v7, v7, v8
	v_mul_f32_e32 v8, v8, v7
	s_waitcnt vmcnt(0) lgkmcnt(0)
	v_mac_f32_e32 v8, s0, v6
	flat_store_dword v[2:3], v8
	flat_load_dword v4, v[4:5]
	v_sqrt_f32_e32 v6, v8
	s_mov_b32 s0, 0x6f800000
	v_cvt_f64_f32_e32 v[6:7], v6
	v_add_f64 v[6:7], v[6:7], s[2:3]
	v_cvt_f32_f64_e32 v6, v[6:7]
	v_cmp_gt_f32_e64 vcc, |v6|, s0
	v_cndmask_b32_e32 v2, 1.0, v9, vcc
	v_mul_f32_e32 v3, v6, v2
	v_rcp_f32_e32 v3, v3
	s_waitcnt vmcnt(0) lgkmcnt(0)
	v_mul_f32_e32 v3, v4, v3
	v_mul_f32_e32 v2, v2, v3
	flat_load_dword v3, v[0:1]
	s_waitcnt vmcnt(0) lgkmcnt(0)
	v_mad_f32 v2, -s1, v2, v3
	flat_store_dword v[0:1], v2
BB8_2:
	s_endpgm
.Lfunc_end8:
	.size	rmsProp, .Lfunc_end8-rmsProp
                                        ; -- End function
	.section	.AMDGPU.csdata
; Kernel info:
; codeLenInByte = 312
; NumSgprs: 14
; NumVgprs: 11
; ScratchSize: 0
; MemoryBound: 0
; FloatMode: 192
; IeeeMode: 1
; LDSByteSize: 0 bytes/workgroup (compile time only)
; SGPRBlocks: 1
; VGPRBlocks: 2
; NumSGPRsForWavesPerEU: 14
; NumVGPRsForWavesPerEU: 11
; Occupancy: 10
; WaveLimiterHint : 1
; COMPUTE_PGM_RSRC2:USER_SGPR: 8
; COMPUTE_PGM_RSRC2:TRAP_HANDLER: 0
; COMPUTE_PGM_RSRC2:TGID_X_EN: 1
; COMPUTE_PGM_RSRC2:TGID_Y_EN: 0
; COMPUTE_PGM_RSRC2:TGID_Z_EN: 0
; COMPUTE_PGM_RSRC2:TIDIG_COMP_CNT: 0
	.text
	.protected	adam            ; -- Begin function adam
	.globl	adam
	.p2align	8
	.type	adam,@function
	.amdgpu_hsa_kernel adam
adam:                                   ; @adam
adam$local:
	.amd_kernel_code_t
		amd_code_version_major = 1
		amd_code_version_minor = 2
		amd_machine_kind = 1
		amd_machine_version_major = 8
		amd_machine_version_minor = 0
		amd_machine_version_stepping = 3
		kernel_code_entry_byte_offset = 256
		kernel_code_prefetch_byte_size = 0
		granulated_workitem_vgpr_count = 2
		granulated_wavefront_sgpr_count = 1
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
		user_sgpr_count = 8
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
		enable_sgpr_dispatch_ptr = 1
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
		workgroup_group_segment_byte_size = 0
		gds_segment_byte_size = 0
		kernarg_segment_byte_size = 104
		workgroup_fbarrier_count = 0
		wavefront_sgpr_count = 14
		workitem_vgpr_count = 12
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
	s_load_dwordx4 s[0:3], s[6:7], 0x20
	s_load_dwordx2 s[10:11], s[6:7], 0x30
	s_load_dword s4, s[4:5], 0x4
	s_waitcnt lgkmcnt(0)
	s_and_b32 s4, s4, 0xffff
	s_mul_i32 s8, s8, s4
	v_add_u32_e32 v0, vcc, s8, v0
	v_add_u32_e32 v0, vcc, s10, v0
	v_cmp_ge_i32_e32 vcc, s3, v0
	s_and_saveexec_b64 s[4:5], vcc
	s_cbranch_execz BB9_2
; %bb.1:
	s_load_dwordx8 s[4:11], s[6:7], 0x0
	v_mov_b32_e32 v1, 0
	v_mov_b32_e32 v2, v0
	v_ashrrev_i64 v[0:1], 30, v[1:2]
	v_sub_f32_e64 v9, 1.0, s0
	s_waitcnt lgkmcnt(0)
	v_mov_b32_e32 v3, s11
	v_add_u32_e32 v2, vcc, s10, v0
	v_addc_u32_e32 v3, vcc, v3, v1, vcc
	v_mov_b32_e32 v5, s7
	v_add_u32_e32 v4, vcc, s6, v0
	v_addc_u32_e32 v5, vcc, v5, v1, vcc
	flat_load_dword v11, v[4:5]
	flat_load_dword v8, v[2:3]
	v_mov_b32_e32 v7, s9
	v_add_u32_e32 v6, vcc, s8, v0
	v_addc_u32_e32 v7, vcc, v7, v1, vcc
	v_sub_f32_e64 v10, 1.0, s1
	v_add_u32_e32 v0, vcc, s4, v0
	s_waitcnt vmcnt(1) lgkmcnt(1)
	v_mul_f32_e32 v9, v9, v11
	s_waitcnt vmcnt(0) lgkmcnt(0)
	v_mac_f32_e32 v9, s0, v8
	flat_store_dword v[2:3], v9
	flat_load_dword v8, v[6:7]
	flat_load_dword v4, v[4:5]
	s_mov_b32 s0, 0xe2308c3a
	s_waitcnt vmcnt(0) lgkmcnt(0)
	v_mul_f32_e32 v5, v10, v4
	v_mul_f32_e32 v9, v4, v5
	v_mac_f32_e32 v9, s1, v8
	flat_store_dword v[6:7], v9
	v_mov_b32_e32 v10, s5
	flat_load_dword v2, v[2:3]
	v_addc_u32_e32 v1, vcc, v10, v1, vcc
	flat_load_dword v3, v[0:1]
	v_sqrt_f32_e32 v4, v9
	s_mov_b32 s1, 0x3e45798e
	v_mov_b32_e32 v8, 0x2f800000
	v_cvt_f64_f32_e32 v[4:5], v4
	v_add_f64 v[4:5], v[4:5], s[0:1]
	s_mov_b32 s0, 0x6f800000
	v_cvt_f32_f64_e32 v4, v[4:5]
	v_cmp_gt_f32_e64 vcc, |v4|, s0
	v_cndmask_b32_e32 v5, 1.0, v8, vcc
	v_mul_f32_e32 v4, v4, v5
	v_rcp_f32_e32 v4, v4
	s_waitcnt vmcnt(1) lgkmcnt(1)
	v_mul_f32_e32 v2, v2, v4
	v_mul_f32_e32 v2, v5, v2
	s_waitcnt vmcnt(0) lgkmcnt(0)
	v_mad_f32 v2, -s2, v2, v3
	flat_store_dword v[0:1], v2
BB9_2:
	s_endpgm
.Lfunc_end9:
	.size	adam, .Lfunc_end9-adam
                                        ; -- End function
	.section	.AMDGPU.csdata
; Kernel info:
; codeLenInByte = 356
; NumSgprs: 14
; NumVgprs: 12
; ScratchSize: 0
; MemoryBound: 0
; FloatMode: 192
; IeeeMode: 1
; LDSByteSize: 0 bytes/workgroup (compile time only)
; SGPRBlocks: 1
; VGPRBlocks: 2
; NumSGPRsForWavesPerEU: 14
; NumVGPRsForWavesPerEU: 12
; Occupancy: 10
; WaveLimiterHint : 1
; COMPUTE_PGM_RSRC2:USER_SGPR: 8
; COMPUTE_PGM_RSRC2:TRAP_HANDLER: 0
; COMPUTE_PGM_RSRC2:TGID_X_EN: 1
; COMPUTE_PGM_RSRC2:TGID_Y_EN: 0
; COMPUTE_PGM_RSRC2:TGID_Z_EN: 0
; COMPUTE_PGM_RSRC2:TIDIG_COMP_CNT: 0
	.text
	.protected	reluForward     ; -- Begin function reluForward
	.globl	reluForward
	.p2align	8
	.type	reluForward,@function
	.amdgpu_hsa_kernel reluForward
reluForward:                            ; @reluForward
reluForward$local:
	.amd_kernel_code_t
		amd_code_version_major = 1
		amd_code_version_minor = 2
		amd_machine_kind = 1
		amd_machine_version_major = 8
		amd_machine_version_minor = 0
		amd_machine_version_stepping = 3
		kernel_code_entry_byte_offset = 256
		kernel_code_prefetch_byte_size = 0
		granulated_workitem_vgpr_count = 1
		granulated_wavefront_sgpr_count = 1
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
		user_sgpr_count = 8
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
		enable_sgpr_dispatch_ptr = 1
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
		workgroup_group_segment_byte_size = 0
		gds_segment_byte_size = 0
		kernarg_segment_byte_size = 80
		workgroup_fbarrier_count = 0
		wavefront_sgpr_count = 11
		workitem_vgpr_count = 5
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
	s_load_dword s0, s[4:5], 0x4
	s_waitcnt lgkmcnt(0)
	s_and_b32 s0, s0, 0xffff
	s_mul_i32 s8, s8, s0
	s_load_dword s2, s[6:7], 0x10
	s_load_dwordx2 s[0:1], s[6:7], 0x18
	v_add_u32_e32 v0, vcc, s8, v0
	s_waitcnt lgkmcnt(0)
	v_add_u32_e32 v0, vcc, s0, v0
	v_cmp_gt_i32_e32 vcc, s2, v0
	s_and_saveexec_b64 s[0:1], vcc
	s_cbranch_execz BB10_2
; %bb.1:
	s_load_dwordx4 s[0:3], s[6:7], 0x0
	v_mov_b32_e32 v1, 0
	v_mov_b32_e32 v2, v0
	v_ashrrev_i64 v[0:1], 30, v[1:2]
	s_waitcnt lgkmcnt(0)
	v_mov_b32_e32 v3, s1
	v_add_u32_e32 v2, vcc, s0, v0
	v_addc_u32_e32 v3, vcc, v3, v1, vcc
	flat_load_dword v2, v[2:3]
	v_mov_b32_e32 v4, s3
	v_add_u32_e32 v0, vcc, s2, v0
	v_addc_u32_e32 v1, vcc, v4, v1, vcc
	s_waitcnt vmcnt(0) lgkmcnt(0)
	v_mul_f32_e32 v2, 1.0, v2
	v_max_f32_e32 v2, 0, v2
	flat_store_dword v[0:1], v2
BB10_2:
	s_endpgm
.Lfunc_end10:
	.size	reluForward, .Lfunc_end10-reluForward
                                        ; -- End function
	.section	.AMDGPU.csdata
; Kernel info:
; codeLenInByte = 148
; NumSgprs: 11
; NumVgprs: 5
; ScratchSize: 0
; MemoryBound: 0
; FloatMode: 192
; IeeeMode: 1
; LDSByteSize: 0 bytes/workgroup (compile time only)
; SGPRBlocks: 1
; VGPRBlocks: 1
; NumSGPRsForWavesPerEU: 11
; NumVGPRsForWavesPerEU: 5
; Occupancy: 10
; WaveLimiterHint : 1
; COMPUTE_PGM_RSRC2:USER_SGPR: 8
; COMPUTE_PGM_RSRC2:TRAP_HANDLER: 0
; COMPUTE_PGM_RSRC2:TGID_X_EN: 1
; COMPUTE_PGM_RSRC2:TGID_Y_EN: 0
; COMPUTE_PGM_RSRC2:TGID_Z_EN: 0
; COMPUTE_PGM_RSRC2:TIDIG_COMP_CNT: 0
	.text
	.protected	reluBackward    ; -- Begin function reluBackward
	.globl	reluBackward
	.p2align	8
	.type	reluBackward,@function
	.amdgpu_hsa_kernel reluBackward
reluBackward:                           ; @reluBackward
reluBackward$local:
	.amd_kernel_code_t
		amd_code_version_major = 1
		amd_code_version_minor = 2
		amd_machine_kind = 1
		amd_machine_version_major = 8
		amd_machine_version_minor = 0
		amd_machine_version_stepping = 3
		kernel_code_entry_byte_offset = 256
		kernel_code_prefetch_byte_size = 0
		granulated_workitem_vgpr_count = 1
		granulated_wavefront_sgpr_count = 1
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
		user_sgpr_count = 8
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
		enable_sgpr_dispatch_ptr = 1
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
		workgroup_group_segment_byte_size = 0
		gds_segment_byte_size = 0
		kernarg_segment_byte_size = 88
		workgroup_fbarrier_count = 0
		wavefront_sgpr_count = 14
		workitem_vgpr_count = 6
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
	s_load_dword s0, s[4:5], 0x4
	s_waitcnt lgkmcnt(0)
	s_and_b32 s0, s0, 0xffff
	s_mul_i32 s8, s8, s0
	s_load_dword s2, s[6:7], 0x18
	s_load_dwordx2 s[0:1], s[6:7], 0x20
	v_add_u32_e32 v0, vcc, s8, v0
	s_waitcnt lgkmcnt(0)
	v_add_u32_e32 v0, vcc, s0, v0
	v_cmp_gt_i32_e32 vcc, s2, v0
	s_and_saveexec_b64 s[0:1], vcc
	s_cbranch_execz BB11_4
; %bb.1:
	s_load_dwordx4 s[8:11], s[6:7], 0x0
	s_load_dwordx4 s[0:3], s[6:7], 0x10
	v_mov_b32_e32 v2, 0
	v_mov_b32_e32 v3, v0
	v_ashrrev_i64 v[3:4], 30, v[2:3]
	s_waitcnt lgkmcnt(0)
	v_mov_b32_e32 v5, s9
	v_add_u32_e32 v3, vcc, s8, v3
	v_addc_u32_e32 v4, vcc, v5, v4, vcc
	flat_load_dword v3, v[3:4]
	v_ashrrev_i32_e32 v1, 31, v0
	s_waitcnt vmcnt(0) lgkmcnt(0)
	v_cmp_lt_f32_e32 vcc, 0, v3
	s_and_saveexec_b64 s[2:3], vcc
	s_cbranch_execz BB11_3
; %bb.2:
	v_lshlrev_b64 v[2:3], 2, v[0:1]
	v_mov_b32_e32 v4, s11
	v_add_u32_e32 v2, vcc, s10, v2
	v_addc_u32_e32 v3, vcc, v4, v3, vcc
	flat_load_dword v2, v[2:3]
BB11_3:
	s_or_b64 exec, exec, s[2:3]
	v_lshlrev_b64 v[0:1], 2, v[0:1]
	v_mov_b32_e32 v3, s1
	v_add_u32_e32 v0, vcc, s0, v0
	v_addc_u32_e32 v1, vcc, v3, v1, vcc
	s_waitcnt vmcnt(0) lgkmcnt(0)
	flat_store_dword v[0:1], v2
BB11_4:
	s_endpgm
.Lfunc_end11:
	.size	reluBackward, .Lfunc_end11-reluBackward
                                        ; -- End function
	.section	.AMDGPU.csdata
; Kernel info:
; codeLenInByte = 208
; NumSgprs: 14
; NumVgprs: 6
; ScratchSize: 0
; MemoryBound: 0
; FloatMode: 192
; IeeeMode: 1
; LDSByteSize: 0 bytes/workgroup (compile time only)
; SGPRBlocks: 1
; VGPRBlocks: 1
; NumSGPRsForWavesPerEU: 14
; NumVGPRsForWavesPerEU: 6
; Occupancy: 10
; WaveLimiterHint : 1
; COMPUTE_PGM_RSRC2:USER_SGPR: 8
; COMPUTE_PGM_RSRC2:TRAP_HANDLER: 0
; COMPUTE_PGM_RSRC2:TGID_X_EN: 1
; COMPUTE_PGM_RSRC2:TGID_Y_EN: 0
; COMPUTE_PGM_RSRC2:TGID_Z_EN: 0
; COMPUTE_PGM_RSRC2:TIDIG_COMP_CNT: 0
	.ident	"clang version 11.0.0 (/src/external/llvm-project/clang 0383ad1cfb0a8e05b0a020e8632400194628b243)"
	.section	".note.GNU-stack"
	.addrsig
	.amd_amdgpu_isa "amdgcn-amd-amdhsa--gfx803"
	.amd_amdgpu_hsa_metadata
---
Version:         [ 1, 0 ]
Kernels:
  - Name:            transpose_tensor
    SymbolName:      'transpose_tensor@kd'
    Language:        OpenCL C
    LanguageVersion: [ 1, 2 ]
    Args:
      - Name:            in
        TypeName:        'float*'
        Size:            8
        Align:           8
        ValueKind:       GlobalBuffer
        ValueType:       F32
        AddrSpaceQual:   Global
        AccQual:         Default
      - Name:            out
        TypeName:        'float*'
        Size:            8
        Align:           8
        ValueKind:       GlobalBuffer
        ValueType:       F32
        AddrSpaceQual:   Global
        AccQual:         Default
      - Name:            in_size
        TypeName:        'int*'
        Size:            8
        Align:           8
        ValueKind:       GlobalBuffer
        ValueType:       I32
        AddrSpaceQual:   Global
        AccQual:         Default
      - Name:            out_size
        TypeName:        'int*'
        Size:            8
        Align:           8
        ValueKind:       GlobalBuffer
        ValueType:       I32
        AddrSpaceQual:   Global
        AccQual:         Default
      - Name:            order
        TypeName:        'int*'
        Size:            8
        Align:           8
        ValueKind:       GlobalBuffer
        ValueType:       I32
        AddrSpaceQual:   Global
        AccQual:         Default
      - Name:            in_index_buf
        TypeName:        'int*'
        Size:            8
        Align:           8
        ValueKind:       GlobalBuffer
        ValueType:       I32
        AddrSpaceQual:   Global
        AccQual:         Default
      - Name:            out_index_buf
        TypeName:        'int*'
        Size:            8
        Align:           8
        ValueKind:       GlobalBuffer
        ValueType:       I32
        AddrSpaceQual:   Global
        AccQual:         Default
      - Name:            dim
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
      KernargSegmentSize: 120
      GroupSegmentFixedSize: 0
      PrivateSegmentFixedSize: 0
      KernargSegmentAlign: 8
      WavefrontSize:   64
      NumSGPRs:        26
      NumVGPRs:        18
      MaxFlatWorkGroupSize: 256
  - Name:            rotate_tensor
    SymbolName:      'rotate_tensor@kd'
    Language:        OpenCL C
    LanguageVersion: [ 1, 2 ]
    Args:
      - Name:            in
        TypeName:        'float*'
        Size:            8
        Align:           8
        ValueKind:       GlobalBuffer
        ValueType:       F32
        AddrSpaceQual:   Global
        AccQual:         Default
      - Name:            out
        TypeName:        'float*'
        Size:            8
        Align:           8
        ValueKind:       GlobalBuffer
        ValueType:       F32
        AddrSpaceQual:   Global
        AccQual:         Default
      - Name:            in_size
        TypeName:        'int*'
        Size:            8
        Align:           8
        ValueKind:       GlobalBuffer
        ValueType:       I32
        AddrSpaceQual:   Global
        AccQual:         Default
      - Name:            out_size
        TypeName:        'int*'
        Size:            8
        Align:           8
        ValueKind:       GlobalBuffer
        ValueType:       I32
        AddrSpaceQual:   Global
        AccQual:         Default
      - Name:            in_index_buf
        TypeName:        'int*'
        Size:            8
        Align:           8
        ValueKind:       GlobalBuffer
        ValueType:       I32
        AddrSpaceQual:   Global
        AccQual:         Default
      - Name:            out_index_buf
        TypeName:        'int*'
        Size:            8
        Align:           8
        ValueKind:       GlobalBuffer
        ValueType:       I32
        AddrSpaceQual:   Global
        AccQual:         Default
      - Name:            dim
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
      KernargSegmentSize: 112
      GroupSegmentFixedSize: 0
      PrivateSegmentFixedSize: 0
      KernargSegmentAlign: 8
      WavefrontSize:   64
      NumSGPRs:        18
      NumVGPRs:        18
      MaxFlatWorkGroupSize: 256
  - Name:            dilate_tensor
    SymbolName:      'dilate_tensor@kd'
    Language:        OpenCL C
    LanguageVersion: [ 1, 2 ]
    Args:
      - Name:            in
        TypeName:        'float*'
        Size:            8
        Align:           8
        ValueKind:       GlobalBuffer
        ValueType:       F32
        AddrSpaceQual:   Global
        AccQual:         Default
      - Name:            out
        TypeName:        'float*'
        Size:            8
        Align:           8
        ValueKind:       GlobalBuffer
        ValueType:       F32
        AddrSpaceQual:   Global
        AccQual:         Default
      - Name:            in_size
        TypeName:        'int*'
        Size:            8
        Align:           8
        ValueKind:       GlobalBuffer
        ValueType:       I32
        AddrSpaceQual:   Global
        AccQual:         Default
      - Name:            out_size
        TypeName:        'int*'
        Size:            8
        Align:           8
        ValueKind:       GlobalBuffer
        ValueType:       I32
        AddrSpaceQual:   Global
        AccQual:         Default
      - Name:            dilate
        TypeName:        'int*'
        Size:            8
        Align:           8
        ValueKind:       GlobalBuffer
        ValueType:       I32
        AddrSpaceQual:   Global
        AccQual:         Default
      - Name:            in_index_buf
        TypeName:        'int*'
        Size:            8
        Align:           8
        ValueKind:       GlobalBuffer
        ValueType:       I32
        AddrSpaceQual:   Global
        AccQual:         Default
      - Name:            out_index_buf
        TypeName:        'int*'
        Size:            8
        Align:           8
        ValueKind:       GlobalBuffer
        ValueType:       I32
        AddrSpaceQual:   Global
        AccQual:         Default
      - Name:            dim
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
      KernargSegmentSize: 120
      GroupSegmentFixedSize: 0
      PrivateSegmentFixedSize: 0
      KernargSegmentAlign: 8
      WavefrontSize:   64
      NumSGPRs:        30
      NumVGPRs:        18
      MaxFlatWorkGroupSize: 256
  - Name:            softmax_exp
    SymbolName:      'softmax_exp@kd'
    Language:        OpenCL C
    LanguageVersion: [ 1, 2 ]
    Args:
      - Name:            input
        TypeName:        'float*'
        Size:            8
        Align:           8
        ValueKind:       GlobalBuffer
        ValueType:       F32
        AddrSpaceQual:   Global
        AccQual:         Default
      - Name:            output
        TypeName:        'float*'
        Size:            8
        Align:           8
        ValueKind:       GlobalBuffer
        ValueType:       F32
        AddrSpaceQual:   Global
        AccQual:         Default
      - Name:            n
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
      KernargSegmentSize: 80
      GroupSegmentFixedSize: 0
      PrivateSegmentFixedSize: 0
      KernargSegmentAlign: 8
      WavefrontSize:   64
      NumSGPRs:        11
      NumVGPRs:        9
      MaxFlatWorkGroupSize: 256
  - Name:            softmax_div
    SymbolName:      'softmax_div@kd'
    Language:        OpenCL C
    LanguageVersion: [ 1, 2 ]
    Args:
      - Name:            exp_input
        TypeName:        'float*'
        Size:            8
        Align:           8
        ValueKind:       GlobalBuffer
        ValueType:       F32
        AddrSpaceQual:   Global
        AccQual:         Default
      - Name:            out
        TypeName:        'float*'
        Size:            8
        Align:           8
        ValueKind:       GlobalBuffer
        ValueType:       F32
        AddrSpaceQual:   Global
        AccQual:         Default
      - Name:            denominator
        TypeName:        'float*'
        Size:            8
        Align:           8
        ValueKind:       GlobalBuffer
        ValueType:       F32
        AddrSpaceQual:   Global
        AccQual:         Default
      - Name:            num_element
        TypeName:        int
        Size:            4
        Align:           4
        ValueKind:       ByValue
        ValueType:       I32
        AccQual:         Default
      - Name:            batch_size
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
      GroupSegmentFixedSize: 0
      PrivateSegmentFixedSize: 0
      KernargSegmentAlign: 8
      WavefrontSize:   64
      NumSGPRs:        14
      NumVGPRs:        10
      MaxFlatWorkGroupSize: 256
  - Name:            sum_one_axis
    SymbolName:      'sum_one_axis@kd'
    Language:        OpenCL C
    LanguageVersion: [ 1, 2 ]
    Args:
      - Name:            in
        TypeName:        'float*'
        Size:            8
        Align:           8
        ValueKind:       GlobalBuffer
        ValueType:       F32
        AddrSpaceQual:   Global
        AccQual:         Default
      - Name:            out
        TypeName:        'float*'
        Size:            8
        Align:           8
        ValueKind:       GlobalBuffer
        ValueType:       F32
        AddrSpaceQual:   Global
        AccQual:         Default
      - Name:            in_size
        TypeName:        'int*'
        Size:            8
        Align:           8
        ValueKind:       GlobalBuffer
        ValueType:       I32
        AddrSpaceQual:   Global
        AccQual:         Default
      - Name:            out_size
        TypeName:        'int*'
        Size:            8
        Align:           8
        ValueKind:       GlobalBuffer
        ValueType:       I32
        AddrSpaceQual:   Global
        AccQual:         Default
      - Name:            in_dim
        TypeName:        int
        Size:            4
        Align:           4
        ValueKind:       ByValue
        ValueType:       I32
        AccQual:         Default
      - Name:            axis
        TypeName:        int
        Size:            4
        Align:           4
        ValueKind:       ByValue
        ValueType:       I32
        AccQual:         Default
      - Name:            in_index_buf
        TypeName:        'int*'
        Size:            8
        Align:           8
        ValueKind:       GlobalBuffer
        ValueType:       I32
        AddrSpaceQual:   Global
        AccQual:         Default
      - Name:            out_index_buf
        TypeName:        'int*'
        Size:            8
        Align:           8
        ValueKind:       GlobalBuffer
        ValueType:       I32
        AddrSpaceQual:   Global
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
      KernargSegmentSize: 112
      GroupSegmentFixedSize: 0
      PrivateSegmentFixedSize: 0
      KernargSegmentAlign: 8
      WavefrontSize:   64
      NumSGPRs:        26
      NumVGPRs:        17
      MaxFlatWorkGroupSize: 256
  - Name:            scaleAdd
    SymbolName:      'scaleAdd@kd'
    Language:        OpenCL C
    LanguageVersion: [ 1, 2 ]
    Args:
      - Name:            out
        TypeName:        'float*'
        Size:            8
        Align:           8
        ValueKind:       GlobalBuffer
        ValueType:       F32
        AddrSpaceQual:   Global
        AccQual:         Default
      - Name:            in1
        TypeName:        'float*'
        Size:            8
        Align:           8
        ValueKind:       GlobalBuffer
        ValueType:       F32
        AddrSpaceQual:   Global
        AccQual:         Default
      - Name:            in2
        TypeName:        'float*'
        Size:            8
        Align:           8
        ValueKind:       GlobalBuffer
        ValueType:       F32
        AddrSpaceQual:   Global
        AccQual:         Default
      - Name:            alpha
        TypeName:        float
        Size:            4
        Align:           4
        ValueKind:       ByValue
        ValueType:       F32
        AccQual:         Default
      - Name:            beta
        TypeName:        float
        Size:            4
        Align:           4
        ValueKind:       ByValue
        ValueType:       F32
        AccQual:         Default
      - Name:            n
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
      KernargSegmentSize: 96
      GroupSegmentFixedSize: 0
      PrivateSegmentFixedSize: 0
      KernargSegmentAlign: 8
      WavefrontSize:   64
      NumSGPRs:        14
      NumVGPRs:        6
      MaxFlatWorkGroupSize: 256
  - Name:            mul
    SymbolName:      'mul@kd'
    Language:        OpenCL C
    LanguageVersion: [ 1, 2 ]
    Args:
      - Name:            out
        TypeName:        'float*'
        Size:            8
        Align:           8
        ValueKind:       GlobalBuffer
        ValueType:       F32
        AddrSpaceQual:   Global
        AccQual:         Default
      - Name:            in1
        TypeName:        'float*'
        Size:            8
        Align:           8
        ValueKind:       GlobalBuffer
        ValueType:       F32
        AddrSpaceQual:   Global
        AccQual:         Default
      - Name:            in2
        TypeName:        'float*'
        Size:            8
        Align:           8
        ValueKind:       GlobalBuffer
        ValueType:       F32
        AddrSpaceQual:   Global
        AccQual:         Default
      - Name:            n
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
      GroupSegmentFixedSize: 0
      PrivateSegmentFixedSize: 0
      KernargSegmentAlign: 8
      WavefrontSize:   64
      NumSGPRs:        11
      NumVGPRs:        6
      MaxFlatWorkGroupSize: 256
  - Name:            rmsProp
    SymbolName:      'rmsProp@kd'
    Language:        OpenCL C
    LanguageVersion: [ 1, 2 ]
    Args:
      - Name:            params
        TypeName:        'float*'
        Size:            8
        Align:           8
        ValueKind:       GlobalBuffer
        ValueType:       F32
        AddrSpaceQual:   Global
        AccQual:         Default
      - Name:            gradients
        TypeName:        'float*'
        Size:            8
        Align:           8
        ValueKind:       GlobalBuffer
        ValueType:       F32
        AddrSpaceQual:   Global
        AccQual:         Default
      - Name:            sHistory
        TypeName:        'float*'
        Size:            8
        Align:           8
        ValueKind:       GlobalBuffer
        ValueType:       F32
        AddrSpaceQual:   Global
        AccQual:         Default
      - Name:            smoothFactor
        TypeName:        float
        Size:            4
        Align:           4
        ValueKind:       ByValue
        ValueType:       F32
        AccQual:         Default
      - Name:            learningRate
        TypeName:        float
        Size:            4
        Align:           4
        ValueKind:       ByValue
        ValueType:       F32
        AccQual:         Default
      - Name:            n
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
      KernargSegmentSize: 96
      GroupSegmentFixedSize: 0
      PrivateSegmentFixedSize: 0
      KernargSegmentAlign: 8
      WavefrontSize:   64
      NumSGPRs:        14
      NumVGPRs:        11
      MaxFlatWorkGroupSize: 256
  - Name:            adam
    SymbolName:      'adam@kd'
    Language:        OpenCL C
    LanguageVersion: [ 1, 2 ]
    Args:
      - Name:            params
        TypeName:        'float*'
        Size:            8
        Align:           8
        ValueKind:       GlobalBuffer
        ValueType:       F32
        AddrSpaceQual:   Global
        AccQual:         Default
      - Name:            gradients
        TypeName:        'float*'
        Size:            8
        Align:           8
        ValueKind:       GlobalBuffer
        ValueType:       F32
        AddrSpaceQual:   Global
        AccQual:         Default
      - Name:            sHistory
        TypeName:        'float*'
        Size:            8
        Align:           8
        ValueKind:       GlobalBuffer
        ValueType:       F32
        AddrSpaceQual:   Global
        AccQual:         Default
      - Name:            vHistory
        TypeName:        'float*'
        Size:            8
        Align:           8
        ValueKind:       GlobalBuffer
        ValueType:       F32
        AddrSpaceQual:   Global
        AccQual:         Default
      - Name:            smoothFactor1
        TypeName:        float
        Size:            4
        Align:           4
        ValueKind:       ByValue
        ValueType:       F32
        AccQual:         Default
      - Name:            smoothFactor2
        TypeName:        float
        Size:            4
        Align:           4
        ValueKind:       ByValue
        ValueType:       F32
        AccQual:         Default
      - Name:            learningRate
        TypeName:        float
        Size:            4
        Align:           4
        ValueKind:       ByValue
        ValueType:       F32
        AccQual:         Default
      - Name:            n
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
      KernargSegmentSize: 104
      GroupSegmentFixedSize: 0
      PrivateSegmentFixedSize: 0
      KernargSegmentAlign: 8
      WavefrontSize:   64
      NumSGPRs:        14
      NumVGPRs:        12
      MaxFlatWorkGroupSize: 256
  - Name:            reluForward
    SymbolName:      'reluForward@kd'
    Language:        OpenCL C
    LanguageVersion: [ 1, 2 ]
    Args:
      - Name:            in
        TypeName:        'float*'
        Size:            8
        Align:           8
        ValueKind:       GlobalBuffer
        ValueType:       F32
        AddrSpaceQual:   Global
        AccQual:         Default
      - Name:            out
        TypeName:        'float*'
        Size:            8
        Align:           8
        ValueKind:       GlobalBuffer
        ValueType:       F32
        AddrSpaceQual:   Global
        AccQual:         Default
      - Name:            count
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
      KernargSegmentSize: 80
      GroupSegmentFixedSize: 0
      PrivateSegmentFixedSize: 0
      KernargSegmentAlign: 8
      WavefrontSize:   64
      NumSGPRs:        11
      NumVGPRs:        5
      MaxFlatWorkGroupSize: 256
  - Name:            reluBackward
    SymbolName:      'reluBackward@kd'
    Language:        OpenCL C
    LanguageVersion: [ 1, 2 ]
    Args:
      - Name:            in
        TypeName:        'float*'
        Size:            8
        Align:           8
        ValueKind:       GlobalBuffer
        ValueType:       F32
        AddrSpaceQual:   Global
        AccQual:         Default
      - Name:            backin
        TypeName:        'float*'
        Size:            8
        Align:           8
        ValueKind:       GlobalBuffer
        ValueType:       F32
        AddrSpaceQual:   Global
        AccQual:         Default
      - Name:            out
        TypeName:        'float*'
        Size:            8
        Align:           8
        ValueKind:       GlobalBuffer
        ValueType:       F32
        AddrSpaceQual:   Global
        AccQual:         Default
      - Name:            count
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
      GroupSegmentFixedSize: 0
      PrivateSegmentFixedSize: 0
      KernargSegmentAlign: 8
      WavefrontSize:   64
      NumSGPRs:        14
      NumVGPRs:        6
      MaxFlatWorkGroupSize: 256
...

	.end_amd_amdgpu_hsa_metadata
