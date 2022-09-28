	.text
	.hsa_code_object_version 2,1
	.hsa_code_object_isa 8,0,3,"AMD","AMDGPU"
	.protected	c_CopySrcToComponents ; -- Begin function c_CopySrcToComponents
	.globl	c_CopySrcToComponents
	.p2align	8
	.type	c_CopySrcToComponents,@function
	.amdgpu_hsa_kernel c_CopySrcToComponents
c_CopySrcToComponents:                  ; @c_CopySrcToComponents
c_CopySrcToComponents$local:
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
		workgroup_group_segment_byte_size = 768
		gds_segment_byte_size = 0
		kernarg_segment_byte_size = 96
		workgroup_fbarrier_count = 0
		wavefront_sgpr_count = 22
		workitem_vgpr_count = 7
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
	s_load_dword s0, s[6:7], 0x20
	s_load_dword s1, s[4:5], 0x4
	v_mul_u32_u24_e32 v1, 3, v0
	s_mov_b32 m0, -1
	s_waitcnt lgkmcnt(0)
	v_mov_b32_e32 v4, s19
	s_and_b32 s1, s1, 0xffff
	s_mul_i32 s1, s8, s1
	s_mul_i32 s2, s1, 3
	v_mad_u32_u24 v2, v0, 3, s2
	v_ashrrev_i32_e32 v3, 31, v2
	v_add_u32_e32 v2, vcc, s18, v2
	v_addc_u32_e32 v3, vcc, v4, v3, vcc
	v_add_u32_e32 v4, vcc, 2, v2
	v_addc_u32_e32 v5, vcc, 0, v3, vcc
	flat_load_ubyte v4, v[4:5]
	flat_load_ushort v2, v[2:3]
	v_add_u32_e32 v0, vcc, s1, v0
	v_cmp_gt_i32_e32 vcc, s0, v0
	s_waitcnt vmcnt(0) lgkmcnt(0)
	v_lshrrev_b32_e32 v3, 8, v2
	ds_write_b8 v1, v2
	ds_write_b8 v1, v4 offset:2
	ds_write_b8 v1, v3 offset:1
	s_waitcnt lgkmcnt(0)
	s_barrier
	s_waitcnt lgkmcnt(0)
	s_and_saveexec_b64 s[0:1], vcc
	s_cbranch_execz BB0_2
; %bb.1:
	ds_read_u8 v4, v1 offset:2
	ds_read_u8 v5, v1 offset:1
	ds_read_u8 v1, v1
	s_movk_i32 s0, 0xff80
	v_mov_b32_e32 v3, s13
	s_waitcnt lgkmcnt(0)
	v_add_u32_e32 v6, vcc, s0, v1
	v_ashrrev_i32_e32 v1, 31, v0
	v_lshlrev_b64 v[0:1], 2, v[0:1]
	v_add_u32_e32 v2, vcc, s12, v0
	v_addc_u32_e32 v3, vcc, v3, v1, vcc
	v_add_u32_e32 v5, vcc, s0, v5
	flat_store_dword v[2:3], v6
	v_mov_b32_e32 v3, s15
	v_add_u32_e32 v2, vcc, s14, v0
	v_addc_u32_e32 v3, vcc, v3, v1, vcc
	flat_store_dword v[2:3], v5
	v_add_u32_e32 v2, vcc, s0, v4
	v_mov_b32_e32 v3, s17
	v_add_u32_e32 v0, vcc, s16, v0
	v_addc_u32_e32 v1, vcc, v3, v1, vcc
	flat_store_dword v[0:1], v2
BB0_2:
	s_endpgm
.Lfunc_end0:
	.size	c_CopySrcToComponents, .Lfunc_end0-c_CopySrcToComponents
                                        ; -- End function
	.section	.AMDGPU.csdata
; Kernel info:
; codeLenInByte = 280
; NumSgprs: 22
; NumVgprs: 7
; ScratchSize: 0
; MemoryBound: 0
; FloatMode: 192
; IeeeMode: 1
; LDSByteSize: 768 bytes/workgroup (compile time only)
; SGPRBlocks: 2
; VGPRBlocks: 1
; NumSGPRsForWavesPerEU: 22
; NumVGPRsForWavesPerEU: 7
; Occupancy: 10
; WaveLimiterHint : 1
; COMPUTE_PGM_RSRC2:USER_SGPR: 8
; COMPUTE_PGM_RSRC2:TRAP_HANDLER: 0
; COMPUTE_PGM_RSRC2:TGID_X_EN: 1
; COMPUTE_PGM_RSRC2:TGID_Y_EN: 0
; COMPUTE_PGM_RSRC2:TGID_Z_EN: 0
; COMPUTE_PGM_RSRC2:TIDIG_COMP_CNT: 0
	.text
	.protected	c_CopySrcToComponent ; -- Begin function c_CopySrcToComponent
	.globl	c_CopySrcToComponent
	.p2align	8
	.type	c_CopySrcToComponent,@function
	.amdgpu_hsa_kernel c_CopySrcToComponent
c_CopySrcToComponent:                   ; @c_CopySrcToComponent
c_CopySrcToComponent$local:
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
		workgroup_group_segment_byte_size = 256
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
	s_load_dwordx4 s[0:3], s[6:7], 0x0
	s_load_dword s6, s[6:7], 0x10
	s_load_dword s4, s[4:5], 0x4
	s_mov_b32 m0, -1
	s_waitcnt lgkmcnt(0)
	v_mov_b32_e32 v4, s3
	s_and_b32 s4, s4, 0xffff
	s_mul_i32 s8, s8, s4
	v_add_u32_e32 v1, vcc, s8, v0
	v_ashrrev_i32_e32 v2, 31, v1
	v_add_u32_e32 v3, vcc, s2, v1
	v_addc_u32_e32 v4, vcc, v4, v2, vcc
	flat_load_ubyte v3, v[3:4]
	v_cmp_gt_i32_e32 vcc, s6, v1
	s_waitcnt vmcnt(0) lgkmcnt(0)
	ds_write_b8 v0, v3
	s_waitcnt lgkmcnt(0)
	s_barrier
	s_waitcnt lgkmcnt(0)
	s_and_saveexec_b64 s[2:3], vcc
	s_cbranch_execz BB1_2
; %bb.1:
	ds_read_u8 v0, v0
	v_lshlrev_b64 v[1:2], 2, v[1:2]
	v_mov_b32_e32 v3, s1
	v_add_u32_e32 v1, vcc, s0, v1
	v_addc_u32_e32 v2, vcc, v3, v2, vcc
	s_waitcnt lgkmcnt(0)
	v_add_u32_e32 v0, vcc, 0xffffff80, v0
	flat_store_dword v[1:2], v0
BB1_2:
	s_endpgm
.Lfunc_end1:
	.size	c_CopySrcToComponent, .Lfunc_end1-c_CopySrcToComponent
                                        ; -- End function
	.section	.AMDGPU.csdata
; Kernel info:
; codeLenInByte = 160
; NumSgprs: 11
; NumVgprs: 5
; ScratchSize: 0
; MemoryBound: 0
; FloatMode: 192
; IeeeMode: 1
; LDSByteSize: 256 bytes/workgroup (compile time only)
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
	.protected	cl_fdwt53Kernel ; -- Begin function cl_fdwt53Kernel
	.globl	cl_fdwt53Kernel
	.p2align	8
	.type	cl_fdwt53Kernel,@function
	.amdgpu_hsa_kernel cl_fdwt53Kernel
cl_fdwt53Kernel:                        ; @cl_fdwt53Kernel
cl_fdwt53Kernel$local:
	.amd_kernel_code_t
		amd_code_version_major = 1
		amd_code_version_minor = 2
		amd_machine_kind = 1
		amd_machine_version_major = 8
		amd_machine_version_minor = 0
		amd_machine_version_stepping = 3
		kernel_code_entry_byte_offset = 256
		kernel_code_prefetch_byte_size = 0
		granulated_workitem_vgpr_count = 10
		granulated_wavefront_sgpr_count = 4
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
		enable_sgpr_workgroup_id_y = 1
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
		workgroup_group_segment_byte_size = 8796
		gds_segment_byte_size = 0
		kernarg_segment_byte_size = 96
		workgroup_fbarrier_count = 0
		wavefront_sgpr_count = 40
		workitem_vgpr_count = 42
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
	s_load_dwordx4 s[8:11], s[4:5], 0x0
	s_load_dwordx4 s[12:15], s[4:5], 0x10
	s_load_dword s0, s[4:5], 0x20
	v_mov_b32_e32 v1, 0
	s_mov_b32 m0, -1
	s_waitcnt lgkmcnt(0)
	v_mov_b32_e32 v2, s15
	v_mov_b32_e32 v3, s0
	ds_write2_b32 v1, v2, v3 offset1:1
	s_mov_b32 s16, s12
	v_mov_b32_e32 v2, 0
	s_branch BB2_2
BB2_1:                                  ;   in Loop: Header=BB2_2 Depth=1
	s_mov_b64 s[2:3], -1
                                        ; implicit-def: $vgpr1
	s_andn2_b64 vcc, exec, s[2:3]
	s_cbranch_vccz BB2_4
BB2_2:                                  ; =>This Inner Loop Header: Depth=1
	v_cmp_ne_u32_e32 vcc, 0x2200, v1
	s_and_b64 vcc, exec, vcc
	ds_write2_b32 v1, v2, v2 offset0:18 offset1:19
	ds_write2_b32 v1, v2, v2 offset0:16 offset1:17
	ds_write2_b32 v1, v2, v2 offset0:20 offset1:21
	s_cbranch_vccz BB2_1
; %bb.3:                                ;   in Loop: Header=BB2_2 Depth=1
	ds_write2_b32 v1, v2, v2 offset0:24 offset1:25
	ds_write2_b32 v1, v2, v2 offset0:22 offset1:23
	ds_write2_b32 v1, v2, v2 offset0:28 offset1:29
	ds_write2_b32 v1, v2, v2 offset0:26 offset1:27
	ds_write2_b32 v1, v2, v2 offset0:30 offset1:31
	v_add_u32_e32 v1, vcc, 64, v1
	s_mov_b64 s[2:3], 0
	s_andn2_b64 vcc, exec, s[2:3]
	s_cbranch_vccnz BB2_2
BB2_4:
	s_lshr_b32 s2, s15, 31
	s_add_i32 s2, s15, s2
	s_ashr_i32 s22, s2, 1
	s_add_i32 s1, s0, 3
	s_add_i32 s2, s22, 2
	s_mul_i32 s18, s2, s1
	v_mov_b32_e32 v3, s1
	s_add_i32 s1, s18, 16
	s_ashr_i32 s3, s1, 31
	s_lshr_b32 s3, s3, 27
	s_add_i32 s3, s1, s3
	s_andn2_b32 s3, s3, 31
	s_sub_i32 s19, s1, s3
	v_mov_b32_e32 v1, 0
	v_mov_b32_e32 v2, s15
	ds_write2_b32 v1, v2, v3 offset0:9 offset1:10
	s_sub_i32 s1, 32, s19
	v_mov_b32_e32 v2, s2
	v_mov_b32_e32 v3, 32
	ds_write2_b32 v1, v2, v3 offset0:11 offset1:12
	s_add_i32 s20, s1, s18
	v_mov_b32_e32 v3, s1
	s_add_i32 s1, s7, 1
	s_mul_i32 s1, s1, s14
	s_mul_i32 s1, s1, s0
	v_mov_b32_e32 v2, s18
	s_add_i32 s1, s1, 1
	s_mul_i32 s21, s7, s14
	ds_write2_b32 v1, v2, v3 offset0:13 offset1:14
	v_mov_b32_e32 v2, s20
	ds_write_b32 v1, v2 offset:60
	v_mov_b32_e32 v2, s2
	s_cmp_ge_i32 s1, s13
	s_mov_b64 s[2:3], -1
	s_mul_i32 s21, s21, s0
	s_mul_i32 s23, s15, s6
	ds_write_b32 v1, v2 offset:8792
	s_cbranch_scc0 BB2_131
; %bb.5:
	v_add_u32_e32 v1, vcc, s23, v0
	v_cmp_ne_u32_e64 s[4:5], s7, 0
	s_sub_i32 s17, 0, s12
	s_and_b64 vcc, exec, s[4:5]
	s_cbranch_vccz BB2_11
; %bb.6:
	v_cmp_le_i32_e32 vcc, s12, v1
                                        ; implicit-def: $vgpr3
                                        ; implicit-def: $sgpr2
	s_and_saveexec_b64 s[0:1], vcc
	s_xor_b64 s[0:1], exec, s[0:1]
; %bb.7:
	s_lshl_b32 s2, s12, 1
	v_sub_u32_e32 v2, vcc, s2, v1
	v_add_u32_e32 v3, vcc, -2, v2
; %bb.8:                                ; %Flow422
	s_or_saveexec_b64 s[0:1], s[0:1]
	v_mov_b32_e32 v4, s2
	s_xor_b64 exec, exec, s[0:1]
; %bb.9:
	v_ashrrev_i32_e32 v2, 31, v1
	v_add_u32_e32 v3, vcc, v1, v2
	s_lshl_b32 s2, s12, 1
	v_xor_b32_e32 v3, v3, v2
	v_mov_b32_e32 v4, s2
; %bb.10:
	s_or_b64 exec, exec, s[0:1]
	s_add_i32 s0, s21, -2
	s_mul_i32 s1, s13, s12
	s_mul_i32 s0, s0, s12
	v_add_u32_e32 v2, vcc, s1, v3
	v_add_u32_e32 v3, vcc, s0, v3
	v_mov_b32_e32 v5, s1
	v_cmp_eq_u32_e32 vcc, s0, v5
	v_cndmask_b32_e32 v4, 0, v4, vcc
	v_sub_u32_e64 v5, s[0:1], v3, v4
	v_mov_b32_e32 v3, s12
	v_mov_b32_e32 v4, s17
	v_cndmask_b32_e32 v3, v3, v4, vcc
	v_add_u32_e32 v4, vcc, v5, v3
	v_sub_u32_e32 v8, vcc, 0, v3
	v_lshlrev_b32_e32 v6, 1, v3
	v_cmp_eq_u32_e32 vcc, v4, v2
	v_cndmask_b32_e32 v6, 0, v6, vcc
	v_cndmask_b32_e32 v3, v3, v8, vcc
	v_sub_u32_e64 v7, s[0:1], v4, v6
	v_add_u32_e32 v4, vcc, v7, v3
	v_sub_u32_e32 v8, vcc, 0, v3
	v_lshlrev_b32_e32 v6, 1, v3
	v_cmp_eq_u32_e32 vcc, v4, v2
	v_cndmask_b32_e32 v6, 0, v6, vcc
	v_sub_u32_e64 v4, s[0:1], v4, v6
	v_ashrrev_i32_e32 v6, 31, v5
	v_lshlrev_b64 v[5:6], 2, v[5:6]
	v_cndmask_b32_e32 v3, v3, v8, vcc
	v_mov_b32_e32 v8, s9
	v_add_u32_e32 v9, vcc, s8, v5
	v_addc_u32_e32 v10, vcc, v8, v6, vcc
	v_ashrrev_i32_e32 v8, 31, v7
	v_lshlrev_b64 v[5:6], 2, v[7:8]
	v_mov_b32_e32 v7, s9
	v_add_u32_e32 v11, vcc, s8, v5
	v_ashrrev_i32_e32 v5, 31, v4
	v_addc_u32_e32 v12, vcc, v7, v6, vcc
	v_lshlrev_b64 v[5:6], 2, v[4:5]
	s_mov_b64 s[0:1], 0
	v_add_u32_e32 v13, vcc, s8, v5
	v_addc_u32_e32 v14, vcc, v7, v6, vcc
	flat_load_dword v7, v[9:10]
	flat_load_dword v6, v[11:12]
	flat_load_dword v5, v[13:14]
	s_and_b64 vcc, exec, s[0:1]
	s_cbranch_vccnz BB2_12
	s_branch BB2_21
BB2_11:
                                        ; implicit-def: $vgpr2
                                        ; implicit-def: $vgpr3
                                        ; implicit-def: $vgpr5
                                        ; implicit-def: $vgpr6
                                        ; implicit-def: $vgpr4
                                        ; implicit-def: $vgpr7
	s_cbranch_execz BB2_21
BB2_12:
	v_cmp_le_i32_e64 s[0:1], s12, v1
                                        ; implicit-def: $vgpr3
                                        ; implicit-def: $sgpr24
	s_and_saveexec_b64 s[2:3], s[0:1]
	s_xor_b64 s[2:3], exec, s[2:3]
; %bb.13:
	s_lshl_b32 s24, s12, 1
	v_sub_u32_e32 v2, vcc, s24, v1
	v_add_u32_e32 v3, vcc, -2, v2
; %bb.14:                               ; %Flow424
	s_or_saveexec_b64 s[2:3], s[2:3]
	v_mov_b32_e32 v2, s24
	s_xor_b64 exec, exec, s[2:3]
; %bb.15:
	v_ashrrev_i32_e32 v2, 31, v1
	v_add_u32_e32 v3, vcc, v1, v2
	s_lshl_b32 s24, s12, 1
	v_xor_b32_e32 v3, v3, v2
	v_mov_b32_e32 v2, s24
; %bb.16:
	s_or_b64 exec, exec, s[2:3]
	s_mul_i32 s24, s13, s12
	s_mul_i32 s25, s21, s12
	v_add_u32_e32 v4, vcc, s24, v3
	v_add_u32_e32 v3, vcc, s25, v3
	s_waitcnt vmcnt(0) lgkmcnt(0)
	v_mov_b32_e32 v5, s24
	v_cmp_eq_u32_e32 vcc, s25, v5
	v_cndmask_b32_e32 v5, 0, v2, vcc
	v_sub_u32_e64 v3, s[2:3], v3, v5
	v_mov_b32_e32 v5, s12
	v_mov_b32_e32 v6, s17
	v_cndmask_b32_e32 v6, v5, v6, vcc
	v_add_u32_e32 v5, vcc, v3, v6
	v_sub_u32_e32 v8, vcc, 0, v6
	v_lshlrev_b32_e32 v7, 1, v6
	v_cmp_eq_u32_e32 vcc, v5, v4
	v_cndmask_b32_e32 v7, 0, v7, vcc
	v_cndmask_b32_e32 v6, v6, v8, vcc
	v_sub_u32_e64 v5, s[2:3], v5, v7
	v_add_u32_e32 v7, vcc, v5, v6
	v_lshlrev_b32_e32 v6, 1, v6
	v_cmp_eq_u32_e32 vcc, v7, v4
	v_cndmask_b32_e32 v4, 0, v6, vcc
	v_sub_u32_e32 v7, vcc, v7, v4
	v_ashrrev_i32_e32 v4, 31, v3
	v_lshlrev_b64 v[3:4], 2, v[3:4]
	v_mov_b32_e32 v6, s9
	v_add_u32_e32 v3, vcc, s8, v3
	v_addc_u32_e32 v4, vcc, v6, v4, vcc
	v_ashrrev_i32_e32 v6, 31, v5
	v_lshlrev_b64 v[5:6], 2, v[5:6]
	v_mov_b32_e32 v8, s9
	v_add_u32_e32 v9, vcc, s8, v5
	v_addc_u32_e32 v10, vcc, v8, v6, vcc
	v_ashrrev_i32_e32 v8, 31, v7
	v_lshlrev_b64 v[5:6], 2, v[7:8]
	v_mov_b32_e32 v8, s9
	v_add_u32_e32 v7, vcc, s8, v5
	v_addc_u32_e32 v8, vcc, v8, v6, vcc
	flat_load_dword v5, v[3:4]
	flat_load_dword v6, v[9:10]
	flat_load_dword v7, v[7:8]
                                        ; implicit-def: $vgpr3
	s_and_saveexec_b64 s[2:3], s[0:1]
	s_xor_b64 s[0:1], exec, s[2:3]
; %bb.17:
	v_sub_u32_e32 v2, vcc, v2, v1
	v_add_u32_e32 v3, vcc, -2, v2
; %bb.18:                               ; %Flow423
	s_or_saveexec_b64 s[0:1], s[0:1]
	s_xor_b64 exec, exec, s[0:1]
; %bb.19:
	v_ashrrev_i32_e32 v2, 31, v1
	v_add_u32_e32 v1, vcc, v1, v2
	v_xor_b32_e32 v3, v1, v2
; %bb.20:
	s_or_b64 exec, exec, s[0:1]
	v_add_u32_e32 v2, vcc, s24, v3
	v_add_u32_e32 v4, vcc, s25, v3
	v_mov_b32_e32 v3, s12
BB2_21:                                 ; %Flow426
	v_add_u32_e32 v8, vcc, 2, v0
	v_and_b32_e32 v1, 1, v8
	v_mul_lo_u32 v9, s20, v1
	v_mov_b32_e32 v1, 0
	v_lshrrev_b32_e32 v8, 1, v8
	v_cmp_gt_u32_e64 s[0:1], 3, v0
	v_add_u32_e32 v15, vcc, v9, v8
	v_mov_b32_e32 v11, v1
	v_mov_b32_e32 v12, v1
	v_mov_b32_e32 v13, v1
	v_mov_b32_e32 v20, v1
	v_mov_b32_e32 v10, v1
	v_mov_b32_e32 v9, v1
	v_mov_b32_e32 v8, v1
	s_and_saveexec_b64 s[24:25], s[0:1]
	s_cbranch_execz BB2_39
; %bb.22:
	v_mov_b32_e32 v8, s15
	v_cmp_eq_u32_e32 vcc, 0, v0
	v_cndmask_b32_e32 v8, -3, v8, vcc
	v_add_u32_e32 v14, vcc, v8, v0
	v_add_u32_e32 v16, vcc, s23, v14
	s_and_b64 vcc, exec, s[4:5]
	s_cbranch_vccz BB2_28
; %bb.23:
	v_cmp_le_i32_e32 vcc, s12, v16
                                        ; implicit-def: $vgpr9
                                        ; implicit-def: $sgpr4
	s_and_saveexec_b64 s[2:3], vcc
	s_xor_b64 s[2:3], exec, s[2:3]
; %bb.24:
	s_lshl_b32 s4, s12, 1
	v_sub_u32_e32 v8, vcc, s4, v16
	v_add_u32_e32 v9, vcc, -2, v8
; %bb.25:                               ; %Flow417
	s_or_saveexec_b64 s[2:3], s[2:3]
	v_mov_b32_e32 v10, s4
	s_xor_b64 exec, exec, s[2:3]
; %bb.26:
	v_ashrrev_i32_e32 v8, 31, v16
	v_add_u32_e32 v9, vcc, v16, v8
	s_lshl_b32 s4, s12, 1
	v_xor_b32_e32 v9, v9, v8
	v_mov_b32_e32 v10, s4
; %bb.27:
	s_or_b64 exec, exec, s[2:3]
	s_add_i32 s2, s21, -2
	s_mul_i32 s3, s13, s12
	s_mul_i32 s2, s2, s12
	v_add_u32_e32 v8, vcc, s3, v9
	v_add_u32_e32 v9, vcc, s2, v9
	v_mov_b32_e32 v11, s3
	v_cmp_eq_u32_e32 vcc, s2, v11
	v_cndmask_b32_e32 v10, 0, v10, vcc
	v_sub_u32_e64 v11, s[2:3], v9, v10
	v_mov_b32_e32 v9, s12
	v_mov_b32_e32 v10, s17
	v_cndmask_b32_e32 v9, v9, v10, vcc
	v_add_u32_e32 v10, vcc, v11, v9
	v_sub_u32_e32 v13, vcc, 0, v9
	v_lshlrev_b32_e32 v12, 1, v9
	v_cmp_eq_u32_e32 vcc, v10, v8
	v_cndmask_b32_e32 v12, 0, v12, vcc
	v_cndmask_b32_e32 v9, v9, v13, vcc
	v_sub_u32_e64 v17, s[2:3], v10, v12
	v_add_u32_e32 v10, vcc, v17, v9
	v_sub_u32_e32 v13, vcc, 0, v9
	v_lshlrev_b32_e32 v12, 1, v9
	v_cmp_eq_u32_e32 vcc, v10, v8
	v_cndmask_b32_e32 v12, 0, v12, vcc
	v_sub_u32_e64 v10, s[2:3], v10, v12
	v_ashrrev_i32_e32 v12, 31, v11
	v_lshlrev_b64 v[11:12], 2, v[11:12]
	v_cndmask_b32_e32 v9, v9, v13, vcc
	v_mov_b32_e32 v13, s9
	v_add_u32_e32 v19, vcc, s8, v11
	v_ashrrev_i32_e32 v18, 31, v17
	v_addc_u32_e32 v20, vcc, v13, v12, vcc
	v_lshlrev_b64 v[11:12], 2, v[17:18]
	s_mov_b64 s[2:3], 0
	v_add_u32_e32 v17, vcc, s8, v11
	v_ashrrev_i32_e32 v11, 31, v10
	v_addc_u32_e32 v18, vcc, v13, v12, vcc
	v_lshlrev_b64 v[11:12], 2, v[10:11]
	v_add_u32_e32 v21, vcc, s8, v11
	v_addc_u32_e32 v22, vcc, v13, v12, vcc
	flat_load_dword v13, v[19:20]
	flat_load_dword v12, v[17:18]
	flat_load_dword v11, v[21:22]
	s_and_b64 vcc, exec, s[2:3]
	s_cbranch_vccnz BB2_29
	s_branch BB2_38
BB2_28:
                                        ; implicit-def: $vgpr8
                                        ; implicit-def: $vgpr9
                                        ; implicit-def: $vgpr10
                                        ; implicit-def: $vgpr13
                                        ; implicit-def: $vgpr12
                                        ; implicit-def: $vgpr11
	s_cbranch_execz BB2_38
BB2_29:
	v_cmp_le_i32_e64 s[2:3], s12, v16
                                        ; implicit-def: $vgpr9
                                        ; implicit-def: $sgpr26
	s_and_saveexec_b64 s[4:5], s[2:3]
	s_xor_b64 s[4:5], exec, s[4:5]
; %bb.30:
	s_lshl_b32 s26, s12, 1
	v_sub_u32_e32 v8, vcc, s26, v16
	v_add_u32_e32 v9, vcc, -2, v8
; %bb.31:                               ; %Flow419
	s_or_saveexec_b64 s[4:5], s[4:5]
	v_mov_b32_e32 v8, s26
	s_xor_b64 exec, exec, s[4:5]
; %bb.32:
	v_ashrrev_i32_e32 v8, 31, v16
	v_add_u32_e32 v9, vcc, v16, v8
	s_lshl_b32 s26, s12, 1
	v_xor_b32_e32 v9, v9, v8
	v_mov_b32_e32 v8, s26
; %bb.33:
	s_or_b64 exec, exec, s[4:5]
	s_mul_i32 s26, s13, s12
	s_mul_i32 s27, s21, s12
	v_add_u32_e32 v10, vcc, s26, v9
	v_add_u32_e32 v9, vcc, s27, v9
	s_waitcnt vmcnt(0) lgkmcnt(0)
	v_mov_b32_e32 v11, s26
	v_cmp_eq_u32_e32 vcc, s27, v11
	v_cndmask_b32_e32 v11, 0, v8, vcc
	v_sub_u32_e64 v9, s[4:5], v9, v11
	v_mov_b32_e32 v11, s12
	v_mov_b32_e32 v12, s17
	v_cndmask_b32_e32 v12, v11, v12, vcc
	v_add_u32_e32 v11, vcc, v9, v12
	v_sub_u32_e32 v17, vcc, 0, v12
	v_lshlrev_b32_e32 v13, 1, v12
	v_cmp_eq_u32_e32 vcc, v11, v10
	v_cndmask_b32_e32 v13, 0, v13, vcc
	v_sub_u32_e64 v11, s[4:5], v11, v13
	v_cndmask_b32_e32 v12, v12, v17, vcc
	v_add_u32_e32 v13, vcc, v11, v12
	v_lshlrev_b32_e32 v12, 1, v12
	v_cmp_eq_u32_e32 vcc, v13, v10
	v_cndmask_b32_e32 v10, 0, v12, vcc
	v_sub_u32_e32 v17, vcc, v13, v10
	v_ashrrev_i32_e32 v10, 31, v9
	v_lshlrev_b64 v[9:10], 2, v[9:10]
	v_mov_b32_e32 v12, s9
	v_add_u32_e32 v9, vcc, s8, v9
	v_addc_u32_e32 v10, vcc, v12, v10, vcc
	v_ashrrev_i32_e32 v12, 31, v11
	v_lshlrev_b64 v[11:12], 2, v[11:12]
	v_mov_b32_e32 v13, s9
	v_add_u32_e32 v19, vcc, s8, v11
	v_ashrrev_i32_e32 v18, 31, v17
	v_addc_u32_e32 v20, vcc, v13, v12, vcc
	v_lshlrev_b64 v[11:12], 2, v[17:18]
	v_add_u32_e32 v17, vcc, s8, v11
	v_addc_u32_e32 v18, vcc, v13, v12, vcc
	flat_load_dword v11, v[9:10]
	flat_load_dword v12, v[19:20]
	flat_load_dword v13, v[17:18]
                                        ; implicit-def: $vgpr9
	s_and_saveexec_b64 s[4:5], s[2:3]
	s_xor_b64 s[2:3], exec, s[4:5]
; %bb.34:
	v_sub_u32_e32 v8, vcc, v8, v16
	v_add_u32_e32 v9, vcc, -2, v8
; %bb.35:                               ; %Flow418
	s_or_saveexec_b64 s[2:3], s[2:3]
	s_xor_b64 exec, exec, s[2:3]
; %bb.36:
	v_ashrrev_i32_e32 v8, 31, v16
	v_add_u32_e32 v9, vcc, v16, v8
	v_xor_b32_e32 v9, v9, v8
; %bb.37:
	s_or_b64 exec, exec, s[2:3]
	v_add_u32_e32 v8, vcc, s26, v9
	v_add_u32_e32 v10, vcc, s27, v9
	v_mov_b32_e32 v9, s12
BB2_38:                                 ; %Flow421
	v_add_u32_e32 v14, vcc, 2, v14
	v_and_b32_e32 v17, 1, v14
	v_mul_lo_u32 v17, v17, s20
	v_lshrrev_b32_e32 v16, 31, v14
	v_add_u32_e32 v14, vcc, v14, v16
	v_ashrrev_i32_e32 v14, 1, v14
	v_add_u32_e32 v20, vcc, v17, v14
BB2_39:
	s_or_b64 exec, exec, s[24:25]
	v_cvt_f32_u32_e32 v14, s22
	v_lshlrev_b64 v[18:19], 1, v[0:1]
	v_cvt_f32_u32_e32 v17, v0
	s_add_i32 s2, s15, -1
	v_rcp_iflag_f32_e32 v19, v14
	s_mul_i32 s3, s15, s6
	v_mov_b32_e32 v16, 0
	v_mov_b32_e32 v23, v16
	v_mul_f32_e32 v19, v17, v19
	v_trunc_f32_e32 v19, v19
	v_cvt_u32_f32_e32 v21, v19
	v_mad_f32 v17, -v19, v14, v17
	v_cmp_ge_f32_e64 vcc, |v17|, v14
	v_mov_b32_e32 v17, v16
	v_addc_u32_e32 v14, vcc, 0, v21, vcc
	v_and_b32_e32 v14, 0x3fffffff, v14
	v_mul_lo_u32 v19, v14, s2
	v_mov_b32_e32 v14, v16
	v_sub_u32_e32 v18, vcc, v18, v19
	v_add_u32_e32 v19, vcc, s3, v18
	v_cmp_gt_i32_e32 vcc, s12, v19
	s_and_saveexec_b64 s[2:3], vcc
	s_cbranch_execz BB2_45
; %bb.40:
	s_lshr_b32 s4, s12, 31
	s_add_i32 s4, s12, s4
	s_ashr_i32 s25, s4, 1
	s_and_b32 s4, s12, 1
	v_lshrrev_b32_e32 v14, 31, v19
	s_add_i32 s17, s25, s4
	s_lshr_b32 s4, s13, 31
	s_add_i32 s4, s13, s4
	v_add_u32_e32 v14, vcc, v19, v14
	v_ashrrev_i32_e32 v16, 1, v14
	v_and_b32_e32 v14, 1, v19
	s_ashr_i32 s4, s4, 1
	s_and_b32 s5, s13, 1
	v_cmp_eq_u32_e32 vcc, 1, v14
	s_add_i32 s24, s4, s5
                                        ; implicit-def: $vgpr19
                                        ; implicit-def: $sgpr28
	s_and_saveexec_b64 s[26:27], vcc
	s_xor_b64 s[26:27], exec, s[26:27]
; %bb.41:
	s_mul_i32 s28, s24, s17
	v_add_u32_e32 v19, vcc, s28, v16
	s_mul_i32 s28, s13, s12
	s_lshr_b32 s29, s28, 31
	s_add_i32 s28, s28, s29
	s_ashr_i32 s28, s28, 1
; %bb.42:                               ; %Flow415
	s_or_saveexec_b64 s[26:27], s[26:27]
	v_mov_b32_e32 v14, s28
	v_mov_b32_e32 v21, s25
	s_xor_b64 exec, exec, s[26:27]
; %bb.43:
	s_mul_i32 s24, s24, s12
	v_mov_b32_e32 v14, s24
	v_mov_b32_e32 v21, s17
	v_mov_b32_e32 v19, v16
; %bb.44:
	s_or_b64 exec, exec, s[26:27]
	v_mul_lo_u32 v16, v14, s5
	v_mul_lo_u32 v22, v21, s4
	s_lshr_b32 s4, s21, 31
	s_add_i32 s4, s21, s4
	v_add_u32_e32 v16, vcc, v16, v19
	s_and_b32 s5, s21, 1
	v_add_u32_e32 v16, vcc, v16, v22
	v_mul_lo_u32 v22, v14, s5
	s_ashr_i32 s4, s4, 1
	v_sub_u32_e32 v17, vcc, v21, v14
	v_mul_lo_u32 v21, v21, s4
	v_add_u32_e32 v19, vcc, v22, v19
	v_add_u32_e32 v23, vcc, v19, v21
BB2_45:                                 ; %Flow416
	s_or_b64 exec, exec, s[2:3]
	v_add_u32_e32 v18, vcc, 2, v18
	v_and_b32_e32 v21, 1, v18
	v_mul_lo_u32 v24, v21, s20
	v_lshrrev_b32_e32 v19, 31, v18
	v_add_u32_e32 v18, vcc, v18, v19
	v_ashrrev_i32_e32 v22, 1, v18
	v_add_u32_e32 v18, vcc, v24, v22
	s_cmp_lt_i32 s14, 1
	s_cbranch_scc1 BB2_130
; %bb.46:
	s_sub_i32 s2, s18, s19
	s_add_i32 s2, s2, 32
	v_mul_lo_u32 v24, s2, v21
	v_lshlrev_b32_e32 v15, 2, v15
	v_add_u32_e32 v19, vcc, 64, v15
	v_lshlrev_b32_e32 v15, 2, v20
	v_add_u32_e32 v20, vcc, 64, v15
	v_add_u32_e32 v21, vcc, v22, v24
	v_lshlrev_b32_e32 v15, 2, v24
	v_lshlrev_b32_e32 v22, 2, v22
	v_add_u32_e32 v22, vcc, v15, v22
	s_mov_b32 s17, 0
	v_mov_b32_e32 v15, 1
	s_mov_b32 s24, 0x7ffffffc
	s_mov_b32 s25, 0x4f800000
	s_movk_i32 s26, 0xff
	s_branch BB2_49
BB2_47:                                 ; %Flow371
                                        ;   in Loop: Header=BB2_49 Depth=1
	v_mov_b32_e32 v23, v24
BB2_48:                                 ; %.loopexit.i
                                        ;   in Loop: Header=BB2_49 Depth=1
	s_add_i32 s17, s17, 1
	s_cmp_eq_u32 s17, s14
	s_waitcnt vmcnt(0) lgkmcnt(0)
	s_barrier
	s_waitcnt lgkmcnt(0)
	s_cbranch_scc1 BB2_130
BB2_49:                                 ; =>This Loop Header: Depth=1
                                        ;     Child Loop BB2_53 Depth 2
                                        ;     Child Loop BB2_57 Depth 2
                                        ;     Child Loop BB2_60 Depth 2
                                        ;     Child Loop BB2_64 Depth 2
                                        ;     Child Loop BB2_67 Depth 2
                                        ;     Child Loop BB2_73 Depth 2
                                        ;     Child Loop BB2_77 Depth 2
                                        ;     Child Loop BB2_80 Depth 2
                                        ;     Child Loop BB2_84 Depth 2
                                        ;     Child Loop BB2_87 Depth 2
                                        ;     Child Loop BB2_91 Depth 2
                                        ;     Child Loop BB2_98 Depth 2
                                        ;     Child Loop BB2_107 Depth 2
                                        ;     Child Loop BB2_127 Depth 2
	s_waitcnt vmcnt(0) lgkmcnt(0)
	ds_write_b32 v19, v7
	v_mov_b32_e32 v7, 0
	ds_read_b32 v24, v7 offset:8792
	s_mov_b64 s[2:3], -1
	s_waitcnt lgkmcnt(0)
	v_lshlrev_b32_e32 v24, 2, v24
	v_add_u32_e32 v24, vcc, v19, v24
	ds_write_b32 v24, v6
	ds_read_b32 v6, v7 offset:8792
	s_waitcnt lgkmcnt(0)
	v_lshlrev_b32_e32 v6, 3, v6
	v_add_u32_e32 v6, vcc, v19, v6
	ds_write_b32 v6, v5
	ds_read_b32 v5, v7 offset:4
                                        ; implicit-def: $vgpr6
	s_waitcnt lgkmcnt(0)
	v_cmp_gt_i32_e32 vcc, 1, v5
	s_and_b64 vcc, exec, vcc
	s_cbranch_vccz BB2_51
; %bb.50:                               ;   in Loop: Header=BB2_49 Depth=1
	v_add_u32_e32 v6, vcc, 2, v5
	s_mov_b64 s[2:3], 0
BB2_51:                                 ; %Flow411
                                        ;   in Loop: Header=BB2_49 Depth=1
	s_andn2_b64 vcc, exec, s[2:3]
	s_cbranch_vccnz BB2_54
; %bb.52:                               ; %.preheader23.i.preheader
                                        ;   in Loop: Header=BB2_49 Depth=1
	s_mov_b32 s4, 2
BB2_53:                                 ; %.preheader23.i
                                        ;   Parent Loop BB2_49 Depth=1
                                        ; =>  This Inner Loop Header: Depth=2
	v_add_u32_e32 v4, vcc, v4, v3
	v_sub_u32_e32 v6, vcc, 0, v3
	v_lshlrev_b32_e32 v5, 1, v3
	v_cmp_eq_u32_e32 vcc, v4, v2
	v_cndmask_b32_e32 v5, 0, v5, vcc
	v_sub_u32_e64 v4, s[2:3], v4, v5
	v_ashrrev_i32_e32 v5, 31, v4
	v_cndmask_b32_e32 v3, v3, v6, vcc
	v_lshlrev_b64 v[5:6], 2, v[4:5]
	v_mov_b32_e32 v7, s9
	v_add_u32_e32 v5, vcc, s8, v5
	v_addc_u32_e32 v6, vcc, v7, v6, vcc
	flat_load_dword v5, v[5:6]
	v_mov_b32_e32 v6, 0
	ds_read_b32 v7, v6 offset:8792
	s_add_i32 s4, s4, 1
	s_waitcnt lgkmcnt(0)
	v_mul_lo_u32 v7, v7, s4
	v_lshlrev_b32_e32 v7, 2, v7
	v_add_u32_e32 v7, vcc, v19, v7
	s_waitcnt vmcnt(0)
	ds_write_b32 v7, v5
	ds_read_b32 v5, v6 offset:4
	s_waitcnt lgkmcnt(0)
	v_add_u32_e32 v6, vcc, 2, v5
	v_cmp_ge_i32_e32 vcc, s4, v6
	s_and_b64 vcc, exec, vcc
	s_cbranch_vccz BB2_53
BB2_54:                                 ; %.loopexit24.i
                                        ;   in Loop: Header=BB2_49 Depth=1
	v_mov_b32_e32 v24, 0
	ds_read_b32 v7, v24 offset:8792
	ds_read_b32 v25, v24 offset:40
	s_waitcnt lgkmcnt(1)
	v_mul_lo_u32 v5, v7, v5
	v_mul_lo_u32 v6, v7, v6
	v_lshlrev_b32_e32 v7, 2, v7
	v_lshlrev_b32_e32 v5, 2, v5
	v_lshlrev_b32_e32 v6, 2, v6
	v_add_u32_e32 v26, vcc, v19, v5
	v_add_u32_e32 v5, vcc, v19, v6
	v_add_u32_e32 v6, vcc, v26, v7
	ds_read_b32 v5, v5
	ds_read_b32 v7, v26
	ds_read_b32 v6, v6
	s_waitcnt lgkmcnt(3)
	v_cmp_gt_i32_e32 vcc, 3, v25
	s_and_b64 vcc, exec, vcc
	s_cbranch_vccnz BB2_68
; %bb.55:                               ;   in Loop: Header=BB2_49 Depth=1
	v_add_u32_e32 v25, vcc, -1, v25
	v_lshrrev_b32_e32 v26, 31, v25
	v_add_u32_e32 v25, vcc, v26, v25
	v_ashrrev_i32_e32 v25, 1, v25
	v_max_i32_e32 v25, 1, v25
	v_add_u32_e32 v26, vcc, -1, v25
	v_cmp_gt_u32_e32 vcc, 3, v26
	s_and_b64 vcc, exec, vcc
	s_cbranch_vccnz BB2_58
; %bb.56:                               ;   in Loop: Header=BB2_49 Depth=1
	v_and_b32_e32 v24, s24, v25
	s_mov_b32 s4, 0
	s_mov_b32 s5, 8
BB2_57:                                 ;   Parent Loop BB2_49 Depth=1
                                        ; =>  This Inner Loop Header: Depth=2
	v_mov_b32_e32 v26, 0
	ds_read_b32 v27, v26 offset:44
	s_add_i32 s27, s5, -8
	s_add_i32 s28, s5, -6
	s_add_i32 s29, s5, -4
	s_add_i32 s30, s5, -2
	s_waitcnt lgkmcnt(0)
	v_mul_lo_u32 v28, s28, v27
	v_mul_lo_u32 v29, v27, s27
	v_lshlrev_b32_e32 v27, 2, v27
	s_add_i32 s4, s4, 4
	v_lshlrev_b32_e32 v28, 2, v28
	v_lshlrev_b32_e32 v29, 2, v29
	v_add_u32_e64 v28, s[2:3], v19, v28
	v_add_u32_e64 v29, s[2:3], v19, v29
	ds_read_b32 v28, v28
	ds_read_b32 v30, v29
	v_add_u32_e64 v27, s[2:3], v29, v27
	ds_read_b32 v29, v27
	v_cmp_ne_u32_e32 vcc, s4, v24
	s_and_b64 vcc, exec, vcc
	s_waitcnt lgkmcnt(1)
	v_add_u32_e64 v28, s[2:3], v30, v28
	v_lshrrev_b32_e32 v30, 31, v28
	v_add_u32_e64 v28, s[2:3], v28, v30
	v_ashrrev_i32_e32 v28, 1, v28
	s_waitcnt lgkmcnt(0)
	v_sub_u32_e64 v28, s[2:3], v29, v28
	ds_write_b32 v27, v28
	ds_read_b32 v27, v26 offset:44
	s_waitcnt lgkmcnt(0)
	v_mul_lo_u32 v28, v27, s29
	v_mul_lo_u32 v29, v27, s28
	v_lshlrev_b32_e32 v27, 2, v27
	v_lshlrev_b32_e32 v28, 2, v28
	v_lshlrev_b32_e32 v29, 2, v29
	v_add_u32_e64 v28, s[2:3], v19, v28
	v_add_u32_e64 v29, s[2:3], v19, v29
	ds_read_b32 v28, v28
	ds_read_b32 v30, v29
	v_add_u32_e64 v27, s[2:3], v29, v27
	ds_read_b32 v29, v27
	s_waitcnt lgkmcnt(1)
	v_add_u32_e64 v28, s[2:3], v30, v28
	v_lshrrev_b32_e32 v30, 31, v28
	v_add_u32_e64 v28, s[2:3], v28, v30
	v_ashrrev_i32_e32 v28, 1, v28
	s_waitcnt lgkmcnt(0)
	v_sub_u32_e64 v28, s[2:3], v29, v28
	ds_write_b32 v27, v28
	ds_read_b32 v27, v26 offset:44
	s_waitcnt lgkmcnt(0)
	v_mul_lo_u32 v28, v27, s30
	v_mul_lo_u32 v29, v27, s29
	v_lshlrev_b32_e32 v27, 2, v27
	v_lshlrev_b32_e32 v28, 2, v28
	v_lshlrev_b32_e32 v29, 2, v29
	v_add_u32_e64 v28, s[2:3], v19, v28
	v_add_u32_e64 v29, s[2:3], v19, v29
	ds_read_b32 v28, v28
	ds_read_b32 v30, v29
	v_add_u32_e64 v27, s[2:3], v29, v27
	ds_read_b32 v29, v27
	s_waitcnt lgkmcnt(1)
	v_add_u32_e64 v28, s[2:3], v30, v28
	v_lshrrev_b32_e32 v30, 31, v28
	v_add_u32_e64 v28, s[2:3], v28, v30
	v_ashrrev_i32_e32 v28, 1, v28
	s_waitcnt lgkmcnt(0)
	v_sub_u32_e64 v28, s[2:3], v29, v28
	ds_write_b32 v27, v28
	ds_read_b32 v26, v26 offset:44
	s_waitcnt lgkmcnt(0)
	v_mul_lo_u32 v27, v26, s5
	v_mul_lo_u32 v28, v26, s30
	v_lshlrev_b32_e32 v26, 2, v26
	s_add_i32 s5, s5, 8
	v_lshlrev_b32_e32 v27, 2, v27
	v_lshlrev_b32_e32 v28, 2, v28
	v_add_u32_e64 v27, s[2:3], v19, v27
	v_add_u32_e64 v28, s[2:3], v19, v28
	ds_read_b32 v27, v27
	ds_read_b32 v29, v28
	v_add_u32_e64 v26, s[2:3], v28, v26
	ds_read_b32 v28, v26
	s_waitcnt lgkmcnt(1)
	v_add_u32_e64 v27, s[2:3], v29, v27
	v_lshrrev_b32_e32 v29, 31, v27
	v_add_u32_e64 v27, s[2:3], v27, v29
	v_ashrrev_i32_e32 v27, 1, v27
	s_waitcnt lgkmcnt(0)
	v_sub_u32_e64 v27, s[2:3], v28, v27
	ds_write_b32 v26, v27
	s_cbranch_vccnz BB2_57
BB2_58:                                 ; %Flow408
                                        ;   in Loop: Header=BB2_49 Depth=1
	v_and_b32_e32 v25, 3, v25
	v_cmp_eq_u32_e32 vcc, 0, v25
	s_and_b64 vcc, exec, vcc
	s_cbranch_vccnz BB2_61
; %bb.59:                               ; %.preheader20.i.preheader
                                        ;   in Loop: Header=BB2_49 Depth=1
	v_lshlrev_b32_e32 v24, 1, v24
BB2_60:                                 ; %.preheader20.i
                                        ;   Parent Loop BB2_49 Depth=1
                                        ; =>  This Inner Loop Header: Depth=2
	v_mov_b32_e32 v26, 0
	ds_read_b32 v26, v26 offset:44
	v_add_u32_e32 v27, vcc, 2, v24
	v_add_u32_e32 v25, vcc, -1, v25
	v_cmp_ne_u32_e32 vcc, 0, v25
	s_and_b64 vcc, exec, vcc
	s_waitcnt lgkmcnt(0)
	v_mul_lo_u32 v28, v27, v26
	v_mul_lo_u32 v29, v26, v24
	v_mov_b32_e32 v24, v27
	v_lshlrev_b32_e32 v26, 2, v26
	v_lshlrev_b32_e32 v28, 2, v28
	v_lshlrev_b32_e32 v27, 2, v29
	v_add_u32_e64 v28, s[2:3], v19, v28
	v_add_u32_e64 v27, s[2:3], v19, v27
	ds_read_b32 v28, v28
	ds_read_b32 v29, v27
	v_add_u32_e64 v26, s[2:3], v27, v26
	ds_read_b32 v27, v26
	s_waitcnt lgkmcnt(1)
	v_add_u32_e64 v28, s[2:3], v29, v28
	v_lshrrev_b32_e32 v29, 31, v28
	v_add_u32_e64 v28, s[2:3], v28, v29
	v_ashrrev_i32_e32 v28, 1, v28
	s_waitcnt lgkmcnt(0)
	v_sub_u32_e64 v27, s[2:3], v27, v28
	ds_write_b32 v26, v27
	s_cbranch_vccnz BB2_60
BB2_61:                                 ; %.loopexit21.i
                                        ;   in Loop: Header=BB2_49 Depth=1
	v_mov_b32_e32 v24, 0
	ds_read_b32 v25, v24 offset:40
	s_waitcnt lgkmcnt(0)
	v_cmp_gt_i32_e32 vcc, 4, v25
	s_and_b64 vcc, exec, vcc
	s_cbranch_vccnz BB2_68
; %bb.62:                               ;   in Loop: Header=BB2_49 Depth=1
	v_lshrrev_b32_e32 v25, 1, v25
	v_max_u32_e32 v26, 2, v25
	v_add_u32_e32 v25, vcc, -1, v26
	v_add_u32_e32 v26, vcc, -2, v26
	v_cmp_gt_u32_e32 vcc, 3, v26
	s_and_b64 vcc, exec, vcc
	s_cbranch_vccnz BB2_65
; %bb.63:                               ;   in Loop: Header=BB2_49 Depth=1
	v_and_b32_e32 v24, -4, v25
	s_mov_b32 s4, 0
	s_mov_b32 s5, 9
BB2_64:                                 ;   Parent Loop BB2_49 Depth=1
                                        ; =>  This Inner Loop Header: Depth=2
	v_mov_b32_e32 v26, 0
	ds_read_b32 v27, v26 offset:44
	s_add_i32 s2, s5, -6
	s_add_i32 s3, s5, -4
	s_add_i32 s27, s5, -2
	s_add_i32 s4, s4, 4
	s_waitcnt lgkmcnt(0)
	v_mul_lo_u32 v28, v27, s2
	v_lshlrev_b32_e32 v30, 1, v27
	v_lshlrev_b32_e32 v27, 2, v27
	v_lshlrev_b32_e32 v29, 2, v28
	v_subrev_u32_e32 v28, vcc, v30, v28
	v_lshlrev_b32_e32 v28, 2, v28
	v_add_u32_e32 v29, vcc, v19, v29
	v_add_u32_e32 v28, vcc, v19, v28
	ds_read_b32 v29, v29
	ds_read_b32 v30, v28
	v_add_u32_e32 v27, vcc, v28, v27
	ds_read_b32 v28, v27
	s_waitcnt lgkmcnt(1)
	v_add_u32_e32 v29, vcc, v29, v30
	v_add_u32_e32 v29, vcc, 2, v29
	v_ashrrev_i32_e32 v30, 31, v29
	v_lshrrev_b32_e32 v30, 30, v30
	v_add_u32_e32 v29, vcc, v29, v30
	v_ashrrev_i32_e32 v29, 2, v29
	s_waitcnt lgkmcnt(0)
	v_add_u32_e32 v28, vcc, v29, v28
	ds_write_b32 v27, v28
	ds_read_b32 v27, v26 offset:44
	s_waitcnt lgkmcnt(0)
	v_mul_lo_u32 v30, v27, s2
	v_mul_lo_u32 v28, v27, s3
	v_lshlrev_b32_e32 v30, 2, v30
	v_lshlrev_b32_e32 v29, 2, v28
	v_add_u32_e32 v29, vcc, v19, v29
	v_add_u32_e32 v30, vcc, v19, v30
	ds_read_b32 v29, v29
	ds_read_b32 v30, v30
	v_subrev_u32_e32 v27, vcc, v27, v28
	v_lshlrev_b32_e32 v27, 2, v27
	v_add_u32_e32 v27, vcc, v19, v27
	ds_read_b32 v28, v27
	s_waitcnt lgkmcnt(1)
	v_add_u32_e32 v29, vcc, v29, v30
	v_add_u32_e32 v29, vcc, 2, v29
	v_ashrrev_i32_e32 v30, 31, v29
	v_lshrrev_b32_e32 v30, 30, v30
	v_add_u32_e32 v29, vcc, v29, v30
	v_ashrrev_i32_e32 v29, 2, v29
	s_waitcnt lgkmcnt(0)
	v_add_u32_e32 v28, vcc, v29, v28
	ds_write_b32 v27, v28
	ds_read_b32 v27, v26 offset:44
	s_waitcnt lgkmcnt(0)
	v_mul_lo_u32 v28, v27, s27
	v_lshlrev_b32_e32 v30, 1, v27
	v_lshlrev_b32_e32 v27, 2, v27
	v_lshlrev_b32_e32 v29, 2, v28
	v_subrev_u32_e32 v28, vcc, v30, v28
	v_lshlrev_b32_e32 v28, 2, v28
	v_add_u32_e32 v29, vcc, v19, v29
	v_add_u32_e32 v28, vcc, v19, v28
	ds_read_b32 v29, v29
	ds_read_b32 v30, v28
	v_add_u32_e32 v27, vcc, v28, v27
	ds_read_b32 v28, v27
	s_waitcnt lgkmcnt(1)
	v_add_u32_e32 v29, vcc, v29, v30
	v_add_u32_e32 v29, vcc, 2, v29
	v_ashrrev_i32_e32 v30, 31, v29
	v_lshrrev_b32_e32 v30, 30, v30
	v_add_u32_e32 v29, vcc, v29, v30
	v_ashrrev_i32_e32 v29, 2, v29
	s_waitcnt lgkmcnt(0)
	v_add_u32_e32 v28, vcc, v29, v28
	ds_write_b32 v27, v28
	ds_read_b32 v26, v26 offset:44
	s_waitcnt lgkmcnt(0)
	v_mul_lo_u32 v27, v26, s5
	v_mul_lo_u32 v29, v26, s27
	s_add_i32 s5, s5, 8
	v_lshlrev_b32_e32 v28, 2, v27
	v_lshlrev_b32_e32 v29, 2, v29
	v_add_u32_e32 v28, vcc, v19, v28
	v_add_u32_e32 v29, vcc, v19, v29
	ds_read_b32 v28, v28
	ds_read_b32 v29, v29
	v_subrev_u32_e32 v26, vcc, v26, v27
	v_lshlrev_b32_e32 v26, 2, v26
	v_add_u32_e64 v26, s[2:3], v19, v26
	ds_read_b32 v27, v26
	s_waitcnt lgkmcnt(1)
	v_add_u32_e64 v28, s[2:3], v28, v29
	v_add_u32_e64 v28, s[2:3], 2, v28
	v_ashrrev_i32_e32 v29, 31, v28
	v_lshrrev_b32_e32 v29, 30, v29
	v_add_u32_e64 v28, s[2:3], v28, v29
	v_cmp_eq_u32_e32 vcc, s4, v24
	v_ashrrev_i32_e32 v28, 2, v28
	s_and_b64 vcc, exec, vcc
	s_waitcnt lgkmcnt(0)
	v_add_u32_e64 v27, s[2:3], v28, v27
	ds_write_b32 v26, v27
	s_cbranch_vccz BB2_64
BB2_65:                                 ; %Flow403
                                        ;   in Loop: Header=BB2_49 Depth=1
	v_and_b32_e32 v25, 3, v25
	v_cmp_eq_u32_e32 vcc, 0, v25
	s_and_b64 vcc, exec, vcc
	s_cbranch_vccnz BB2_68
; %bb.66:                               ; %.preheader17.i.preheader
                                        ;   in Loop: Header=BB2_49 Depth=1
	v_lshlrev_b32_e32 v24, 1, v24
	v_add_u32_e32 v24, vcc, 3, v24
BB2_67:                                 ; %.preheader17.i
                                        ;   Parent Loop BB2_49 Depth=1
                                        ; =>  This Inner Loop Header: Depth=2
	v_mov_b32_e32 v26, 0
	ds_read_b32 v26, v26 offset:44
	v_add_u32_e32 v25, vcc, -1, v25
	s_waitcnt lgkmcnt(0)
	v_mul_lo_u32 v27, v26, v24
	v_lshlrev_b32_e32 v29, 1, v26
	v_lshlrev_b32_e32 v26, 2, v26
	v_add_u32_e32 v24, vcc, 2, v24
	v_lshlrev_b32_e32 v28, 2, v27
	v_subrev_u32_e32 v27, vcc, v29, v27
	v_lshlrev_b32_e32 v27, 2, v27
	v_add_u32_e32 v28, vcc, v19, v28
	v_add_u32_e64 v27, s[2:3], v19, v27
	ds_read_b32 v28, v28
	ds_read_b32 v29, v27
	v_add_u32_e64 v26, s[2:3], v27, v26
	ds_read_b32 v27, v26
	v_cmp_ne_u32_e32 vcc, 0, v25
	s_and_b64 vcc, exec, vcc
	s_waitcnt lgkmcnt(1)
	v_add_u32_e64 v28, s[2:3], v28, v29
	v_add_u32_e64 v28, s[2:3], 2, v28
	v_ashrrev_i32_e32 v29, 31, v28
	v_lshrrev_b32_e32 v29, 30, v29
	v_add_u32_e64 v28, s[2:3], v28, v29
	v_ashrrev_i32_e32 v28, 2, v28
	s_waitcnt lgkmcnt(0)
	v_add_u32_e64 v27, s[2:3], v28, v27
	ds_write_b32 v26, v27
	s_cbranch_vccnz BB2_67
BB2_68:                                 ; %.loopexit18.i
                                        ;   in Loop: Header=BB2_49 Depth=1
	s_and_saveexec_b64 s[4:5], s[0:1]
	s_cbranch_execz BB2_88
; %bb.69:                               ;   in Loop: Header=BB2_49 Depth=1
	ds_write_b32 v20, v13
	v_mov_b32_e32 v13, 0
	ds_read_b32 v24, v13 offset:8792
	s_waitcnt lgkmcnt(0)
	v_lshlrev_b32_e32 v24, 2, v24
	v_add_u32_e32 v24, vcc, v20, v24
	ds_write_b32 v24, v12
	ds_read_b32 v12, v13 offset:8792
	s_waitcnt lgkmcnt(0)
	v_lshlrev_b32_e32 v12, 3, v12
	v_add_u32_e32 v12, vcc, v20, v12
	ds_write_b32 v12, v11
	ds_read_b32 v11, v13 offset:4
	s_waitcnt lgkmcnt(0)
	v_cmp_gt_i32_e32 vcc, 1, v11
	s_and_b64 vcc, exec, vcc
	s_cbranch_vccz BB2_71
; %bb.70:                               ;   in Loop: Header=BB2_49 Depth=1
	v_add_u32_e32 v12, vcc, 2, v11
	s_mov_b64 s[2:3], 0
	s_andn2_b64 vcc, exec, s[2:3]
	s_cbranch_vccz BB2_72
	s_branch BB2_74
BB2_71:                                 ;   in Loop: Header=BB2_49 Depth=1
	s_mov_b64 s[2:3], -1
                                        ; implicit-def: $vgpr12
	s_andn2_b64 vcc, exec, s[2:3]
	s_cbranch_vccnz BB2_74
BB2_72:                                 ; %.preheader15.i.preheader
                                        ;   in Loop: Header=BB2_49 Depth=1
	s_mov_b32 s27, 2
BB2_73:                                 ; %.preheader15.i
                                        ;   Parent Loop BB2_49 Depth=1
                                        ; =>  This Inner Loop Header: Depth=2
	v_add_u32_e32 v10, vcc, v10, v9
	v_sub_u32_e32 v12, vcc, 0, v9
	v_lshlrev_b32_e32 v11, 1, v9
	v_cmp_eq_u32_e32 vcc, v10, v8
	v_cndmask_b32_e32 v11, 0, v11, vcc
	v_sub_u32_e64 v10, s[2:3], v10, v11
	v_ashrrev_i32_e32 v11, 31, v10
	v_cndmask_b32_e32 v9, v9, v12, vcc
	v_lshlrev_b64 v[11:12], 2, v[10:11]
	v_mov_b32_e32 v13, s9
	v_add_u32_e32 v11, vcc, s8, v11
	v_addc_u32_e32 v12, vcc, v13, v12, vcc
	flat_load_dword v11, v[11:12]
	v_mov_b32_e32 v12, 0
	ds_read_b32 v13, v12 offset:8792
	s_add_i32 s27, s27, 1
	s_waitcnt lgkmcnt(0)
	v_mul_lo_u32 v13, v13, s27
	v_lshlrev_b32_e32 v13, 2, v13
	v_add_u32_e32 v13, vcc, v20, v13
	s_waitcnt vmcnt(0)
	ds_write_b32 v13, v11
	ds_read_b32 v11, v12 offset:4
	s_waitcnt lgkmcnt(0)
	v_add_u32_e32 v12, vcc, 2, v11
	v_cmp_ge_i32_e32 vcc, s27, v12
	s_and_b64 vcc, exec, vcc
	s_cbranch_vccz BB2_73
BB2_74:                                 ; %.loopexit16.i
                                        ;   in Loop: Header=BB2_49 Depth=1
	v_mov_b32_e32 v24, 0
	ds_read_b32 v13, v24 offset:8792
	ds_read_b32 v25, v24 offset:40
	s_waitcnt lgkmcnt(1)
	v_mul_lo_u32 v11, v13, v11
	v_mul_lo_u32 v12, v13, v12
	v_lshlrev_b32_e32 v13, 2, v13
	v_lshlrev_b32_e32 v11, 2, v11
	v_lshlrev_b32_e32 v12, 2, v12
	v_add_u32_e32 v26, vcc, v20, v11
	v_add_u32_e32 v11, vcc, v20, v12
	v_add_u32_e32 v12, vcc, v26, v13
	ds_read_b32 v11, v11
	ds_read_b32 v13, v26
	ds_read_b32 v12, v12
	s_waitcnt lgkmcnt(3)
	v_cmp_gt_i32_e32 vcc, 3, v25
	s_and_b64 vcc, exec, vcc
	s_cbranch_vccnz BB2_88
; %bb.75:                               ;   in Loop: Header=BB2_49 Depth=1
	v_add_u32_e32 v25, vcc, -1, v25
	v_lshrrev_b32_e32 v26, 31, v25
	v_add_u32_e32 v25, vcc, v26, v25
	v_ashrrev_i32_e32 v25, 1, v25
	v_max_i32_e32 v25, 1, v25
	v_add_u32_e32 v26, vcc, -1, v25
	v_cmp_gt_u32_e32 vcc, 3, v26
	s_and_b64 vcc, exec, vcc
	s_cbranch_vccnz BB2_78
; %bb.76:                               ;   in Loop: Header=BB2_49 Depth=1
	v_and_b32_e32 v24, s24, v25
	s_mov_b32 s27, 0
	s_mov_b32 s28, 8
BB2_77:                                 ;   Parent Loop BB2_49 Depth=1
                                        ; =>  This Inner Loop Header: Depth=2
	v_mov_b32_e32 v26, 0
	ds_read_b32 v27, v26 offset:44
	s_add_i32 s29, s28, -8
	s_add_i32 s30, s28, -6
	s_add_i32 s31, s28, -4
	s_add_i32 s33, s28, -2
	s_waitcnt lgkmcnt(0)
	v_mul_lo_u32 v28, s30, v27
	v_mul_lo_u32 v29, v27, s29
	v_lshlrev_b32_e32 v27, 2, v27
	s_add_i32 s27, s27, 4
	v_lshlrev_b32_e32 v28, 2, v28
	v_lshlrev_b32_e32 v29, 2, v29
	v_add_u32_e64 v28, s[2:3], v20, v28
	v_add_u32_e64 v29, s[2:3], v20, v29
	ds_read_b32 v28, v28
	ds_read_b32 v30, v29
	v_add_u32_e64 v27, s[2:3], v29, v27
	ds_read_b32 v29, v27
	v_cmp_ne_u32_e32 vcc, s27, v24
	s_and_b64 vcc, exec, vcc
	s_waitcnt lgkmcnt(1)
	v_add_u32_e64 v28, s[2:3], v30, v28
	v_lshrrev_b32_e32 v30, 31, v28
	v_add_u32_e64 v28, s[2:3], v28, v30
	v_ashrrev_i32_e32 v28, 1, v28
	s_waitcnt lgkmcnt(0)
	v_sub_u32_e64 v28, s[2:3], v29, v28
	ds_write_b32 v27, v28
	ds_read_b32 v27, v26 offset:44
	s_waitcnt lgkmcnt(0)
	v_mul_lo_u32 v28, v27, s31
	v_mul_lo_u32 v29, v27, s30
	v_lshlrev_b32_e32 v27, 2, v27
	v_lshlrev_b32_e32 v28, 2, v28
	v_lshlrev_b32_e32 v29, 2, v29
	v_add_u32_e64 v28, s[2:3], v20, v28
	v_add_u32_e64 v29, s[2:3], v20, v29
	ds_read_b32 v28, v28
	ds_read_b32 v30, v29
	v_add_u32_e64 v27, s[2:3], v29, v27
	ds_read_b32 v29, v27
	s_waitcnt lgkmcnt(1)
	v_add_u32_e64 v28, s[2:3], v30, v28
	v_lshrrev_b32_e32 v30, 31, v28
	v_add_u32_e64 v28, s[2:3], v28, v30
	v_ashrrev_i32_e32 v28, 1, v28
	s_waitcnt lgkmcnt(0)
	v_sub_u32_e64 v28, s[2:3], v29, v28
	ds_write_b32 v27, v28
	ds_read_b32 v27, v26 offset:44
	s_waitcnt lgkmcnt(0)
	v_mul_lo_u32 v28, v27, s33
	v_mul_lo_u32 v29, v27, s31
	v_lshlrev_b32_e32 v27, 2, v27
	v_lshlrev_b32_e32 v28, 2, v28
	v_lshlrev_b32_e32 v29, 2, v29
	v_add_u32_e64 v28, s[2:3], v20, v28
	v_add_u32_e64 v29, s[2:3], v20, v29
	ds_read_b32 v28, v28
	ds_read_b32 v30, v29
	v_add_u32_e64 v27, s[2:3], v29, v27
	ds_read_b32 v29, v27
	s_waitcnt lgkmcnt(1)
	v_add_u32_e64 v28, s[2:3], v30, v28
	v_lshrrev_b32_e32 v30, 31, v28
	v_add_u32_e64 v28, s[2:3], v28, v30
	v_ashrrev_i32_e32 v28, 1, v28
	s_waitcnt lgkmcnt(0)
	v_sub_u32_e64 v28, s[2:3], v29, v28
	ds_write_b32 v27, v28
	ds_read_b32 v26, v26 offset:44
	s_waitcnt lgkmcnt(0)
	v_mul_lo_u32 v27, v26, s28
	v_mul_lo_u32 v28, v26, s33
	v_lshlrev_b32_e32 v26, 2, v26
	s_add_i32 s28, s28, 8
	v_lshlrev_b32_e32 v27, 2, v27
	v_lshlrev_b32_e32 v28, 2, v28
	v_add_u32_e64 v27, s[2:3], v20, v27
	v_add_u32_e64 v28, s[2:3], v20, v28
	ds_read_b32 v27, v27
	ds_read_b32 v29, v28
	v_add_u32_e64 v26, s[2:3], v28, v26
	ds_read_b32 v28, v26
	s_waitcnt lgkmcnt(1)
	v_add_u32_e64 v27, s[2:3], v29, v27
	v_lshrrev_b32_e32 v29, 31, v27
	v_add_u32_e64 v27, s[2:3], v27, v29
	v_ashrrev_i32_e32 v27, 1, v27
	s_waitcnt lgkmcnt(0)
	v_sub_u32_e64 v27, s[2:3], v28, v27
	ds_write_b32 v26, v27
	s_cbranch_vccnz BB2_77
BB2_78:                                 ; %Flow394
                                        ;   in Loop: Header=BB2_49 Depth=1
	v_and_b32_e32 v25, 3, v25
	v_cmp_eq_u32_e32 vcc, 0, v25
	s_and_b64 vcc, exec, vcc
	s_cbranch_vccnz BB2_81
; %bb.79:                               ; %.preheader12.i.preheader
                                        ;   in Loop: Header=BB2_49 Depth=1
	v_lshlrev_b32_e32 v24, 1, v24
BB2_80:                                 ; %.preheader12.i
                                        ;   Parent Loop BB2_49 Depth=1
                                        ; =>  This Inner Loop Header: Depth=2
	v_mov_b32_e32 v26, 0
	ds_read_b32 v26, v26 offset:44
	v_add_u32_e32 v27, vcc, 2, v24
	v_add_u32_e32 v25, vcc, -1, v25
	v_cmp_ne_u32_e32 vcc, 0, v25
	s_and_b64 vcc, exec, vcc
	s_waitcnt lgkmcnt(0)
	v_mul_lo_u32 v28, v27, v26
	v_mul_lo_u32 v29, v26, v24
	v_mov_b32_e32 v24, v27
	v_lshlrev_b32_e32 v26, 2, v26
	v_lshlrev_b32_e32 v28, 2, v28
	v_lshlrev_b32_e32 v27, 2, v29
	v_add_u32_e64 v28, s[2:3], v20, v28
	v_add_u32_e64 v27, s[2:3], v20, v27
	ds_read_b32 v28, v28
	ds_read_b32 v29, v27
	v_add_u32_e64 v26, s[2:3], v27, v26
	ds_read_b32 v27, v26
	s_waitcnt lgkmcnt(1)
	v_add_u32_e64 v28, s[2:3], v29, v28
	v_lshrrev_b32_e32 v29, 31, v28
	v_add_u32_e64 v28, s[2:3], v28, v29
	v_ashrrev_i32_e32 v28, 1, v28
	s_waitcnt lgkmcnt(0)
	v_sub_u32_e64 v27, s[2:3], v27, v28
	ds_write_b32 v26, v27
	s_cbranch_vccnz BB2_80
BB2_81:                                 ; %.loopexit13.i
                                        ;   in Loop: Header=BB2_49 Depth=1
	v_mov_b32_e32 v24, 0
	ds_read_b32 v25, v24 offset:40
	s_waitcnt lgkmcnt(0)
	v_cmp_gt_i32_e32 vcc, 4, v25
	s_and_b64 vcc, exec, vcc
	s_cbranch_vccnz BB2_88
; %bb.82:                               ;   in Loop: Header=BB2_49 Depth=1
	v_lshrrev_b32_e32 v25, 1, v25
	v_max_u32_e32 v26, 2, v25
	v_add_u32_e32 v25, vcc, -1, v26
	v_add_u32_e32 v26, vcc, -2, v26
	v_cmp_gt_u32_e32 vcc, 3, v26
	s_and_b64 vcc, exec, vcc
	s_cbranch_vccnz BB2_85
; %bb.83:                               ;   in Loop: Header=BB2_49 Depth=1
	v_and_b32_e32 v24, -4, v25
	s_mov_b32 s27, 0
	s_mov_b32 s28, 9
BB2_84:                                 ;   Parent Loop BB2_49 Depth=1
                                        ; =>  This Inner Loop Header: Depth=2
	v_mov_b32_e32 v26, 0
	ds_read_b32 v27, v26 offset:44
	s_add_i32 s2, s28, -6
	s_add_i32 s3, s28, -4
	s_add_i32 s29, s28, -2
	s_add_i32 s27, s27, 4
	s_waitcnt lgkmcnt(0)
	v_mul_lo_u32 v28, v27, s2
	v_lshlrev_b32_e32 v30, 1, v27
	v_lshlrev_b32_e32 v27, 2, v27
	v_lshlrev_b32_e32 v29, 2, v28
	v_subrev_u32_e32 v28, vcc, v30, v28
	v_lshlrev_b32_e32 v28, 2, v28
	v_add_u32_e32 v29, vcc, v20, v29
	v_add_u32_e32 v28, vcc, v20, v28
	ds_read_b32 v29, v29
	ds_read_b32 v30, v28
	v_add_u32_e32 v27, vcc, v28, v27
	ds_read_b32 v28, v27
	s_waitcnt lgkmcnt(1)
	v_add_u32_e32 v29, vcc, v29, v30
	v_add_u32_e32 v29, vcc, 2, v29
	v_ashrrev_i32_e32 v30, 31, v29
	v_lshrrev_b32_e32 v30, 30, v30
	v_add_u32_e32 v29, vcc, v29, v30
	v_ashrrev_i32_e32 v29, 2, v29
	s_waitcnt lgkmcnt(0)
	v_add_u32_e32 v28, vcc, v29, v28
	ds_write_b32 v27, v28
	ds_read_b32 v27, v26 offset:44
	s_waitcnt lgkmcnt(0)
	v_mul_lo_u32 v30, v27, s2
	v_mul_lo_u32 v28, v27, s3
	v_lshlrev_b32_e32 v30, 2, v30
	v_lshlrev_b32_e32 v29, 2, v28
	v_add_u32_e32 v29, vcc, v20, v29
	v_add_u32_e32 v30, vcc, v20, v30
	ds_read_b32 v29, v29
	ds_read_b32 v30, v30
	v_subrev_u32_e32 v27, vcc, v27, v28
	v_lshlrev_b32_e32 v27, 2, v27
	v_add_u32_e32 v27, vcc, v20, v27
	ds_read_b32 v28, v27
	s_waitcnt lgkmcnt(1)
	v_add_u32_e32 v29, vcc, v29, v30
	v_add_u32_e32 v29, vcc, 2, v29
	v_ashrrev_i32_e32 v30, 31, v29
	v_lshrrev_b32_e32 v30, 30, v30
	v_add_u32_e32 v29, vcc, v29, v30
	v_ashrrev_i32_e32 v29, 2, v29
	s_waitcnt lgkmcnt(0)
	v_add_u32_e32 v28, vcc, v29, v28
	ds_write_b32 v27, v28
	ds_read_b32 v27, v26 offset:44
	s_waitcnt lgkmcnt(0)
	v_mul_lo_u32 v28, v27, s29
	v_lshlrev_b32_e32 v30, 1, v27
	v_lshlrev_b32_e32 v27, 2, v27
	v_lshlrev_b32_e32 v29, 2, v28
	v_subrev_u32_e32 v28, vcc, v30, v28
	v_lshlrev_b32_e32 v28, 2, v28
	v_add_u32_e32 v29, vcc, v20, v29
	v_add_u32_e32 v28, vcc, v20, v28
	ds_read_b32 v29, v29
	ds_read_b32 v30, v28
	v_add_u32_e32 v27, vcc, v28, v27
	ds_read_b32 v28, v27
	s_waitcnt lgkmcnt(1)
	v_add_u32_e32 v29, vcc, v29, v30
	v_add_u32_e32 v29, vcc, 2, v29
	v_ashrrev_i32_e32 v30, 31, v29
	v_lshrrev_b32_e32 v30, 30, v30
	v_add_u32_e32 v29, vcc, v29, v30
	v_ashrrev_i32_e32 v29, 2, v29
	s_waitcnt lgkmcnt(0)
	v_add_u32_e32 v28, vcc, v29, v28
	ds_write_b32 v27, v28
	ds_read_b32 v26, v26 offset:44
	s_waitcnt lgkmcnt(0)
	v_mul_lo_u32 v27, v26, s28
	v_mul_lo_u32 v29, v26, s29
	s_add_i32 s28, s28, 8
	v_lshlrev_b32_e32 v28, 2, v27
	v_lshlrev_b32_e32 v29, 2, v29
	v_add_u32_e32 v28, vcc, v20, v28
	v_add_u32_e32 v29, vcc, v20, v29
	ds_read_b32 v28, v28
	ds_read_b32 v29, v29
	v_subrev_u32_e32 v26, vcc, v26, v27
	v_lshlrev_b32_e32 v26, 2, v26
	v_add_u32_e64 v26, s[2:3], v20, v26
	ds_read_b32 v27, v26
	s_waitcnt lgkmcnt(1)
	v_add_u32_e64 v28, s[2:3], v28, v29
	v_add_u32_e64 v28, s[2:3], 2, v28
	v_ashrrev_i32_e32 v29, 31, v28
	v_lshrrev_b32_e32 v29, 30, v29
	v_add_u32_e64 v28, s[2:3], v28, v29
	v_cmp_eq_u32_e32 vcc, s27, v24
	v_ashrrev_i32_e32 v28, 2, v28
	s_and_b64 vcc, exec, vcc
	s_waitcnt lgkmcnt(0)
	v_add_u32_e64 v27, s[2:3], v28, v27
	ds_write_b32 v26, v27
	s_cbranch_vccz BB2_84
BB2_85:                                 ; %Flow389
                                        ;   in Loop: Header=BB2_49 Depth=1
	v_and_b32_e32 v25, 3, v25
	v_cmp_eq_u32_e32 vcc, 0, v25
	s_and_b64 vcc, exec, vcc
	s_cbranch_vccnz BB2_88
; %bb.86:                               ; %.preheader9.i.preheader
                                        ;   in Loop: Header=BB2_49 Depth=1
	v_lshlrev_b32_e32 v24, 1, v24
	v_add_u32_e32 v24, vcc, 3, v24
BB2_87:                                 ; %.preheader9.i
                                        ;   Parent Loop BB2_49 Depth=1
                                        ; =>  This Inner Loop Header: Depth=2
	v_mov_b32_e32 v26, 0
	ds_read_b32 v26, v26 offset:44
	v_add_u32_e32 v25, vcc, -1, v25
	s_waitcnt lgkmcnt(0)
	v_mul_lo_u32 v27, v26, v24
	v_lshlrev_b32_e32 v29, 1, v26
	v_lshlrev_b32_e32 v26, 2, v26
	v_add_u32_e32 v24, vcc, 2, v24
	v_lshlrev_b32_e32 v28, 2, v27
	v_subrev_u32_e32 v27, vcc, v29, v27
	v_lshlrev_b32_e32 v27, 2, v27
	v_add_u32_e32 v28, vcc, v20, v28
	v_add_u32_e64 v27, s[2:3], v20, v27
	ds_read_b32 v28, v28
	ds_read_b32 v29, v27
	v_add_u32_e64 v26, s[2:3], v27, v26
	ds_read_b32 v27, v26
	v_cmp_ne_u32_e32 vcc, 0, v25
	s_and_b64 vcc, exec, vcc
	s_waitcnt lgkmcnt(1)
	v_add_u32_e64 v28, s[2:3], v28, v29
	v_add_u32_e64 v28, s[2:3], 2, v28
	v_ashrrev_i32_e32 v29, 31, v28
	v_lshrrev_b32_e32 v29, 30, v29
	v_add_u32_e64 v28, s[2:3], v28, v29
	v_ashrrev_i32_e32 v28, 2, v28
	s_waitcnt lgkmcnt(0)
	v_add_u32_e64 v27, s[2:3], v28, v27
	ds_write_b32 v26, v27
	s_cbranch_vccnz BB2_87
BB2_88:                                 ; %Flow399
                                        ;   in Loop: Header=BB2_49 Depth=1
	s_or_b64 exec, exec, s[4:5]
	v_mov_b32_e32 v26, 0
	s_waitcnt lgkmcnt(0)
	s_barrier
	s_waitcnt lgkmcnt(0)
	ds_read2_b32 v[24:25], v26 offset0:1 offset1:9
	s_waitcnt lgkmcnt(0)
	v_ashrrev_i32_e32 v29, 31, v25
	v_add_u32_e32 v27, vcc, v29, v25
	v_xor_b32_e32 v30, v27, v29
	v_cvt_f32_u32_e32 v27, v30
	v_rcp_iflag_f32_e32 v31, v27
	ds_read2_b32 v[27:28], v26 offset0:11 offset1:15
	s_waitcnt lgkmcnt(0)
	v_mul_lo_u32 v26, v27, v24
	v_mul_f32_e32 v24, s25, v31
	v_cvt_u32_f32_e32 v31, v24
	v_lshlrev_b32_e32 v24, 1, v27
	v_add_u32_e32 v27, vcc, -1, v26
	v_ashrrev_i32_e32 v33, 31, v27
	v_mul_lo_u32 v26, v31, v30
	v_mul_hi_u32 v32, v31, v30
	v_xor_b32_e32 v29, v33, v29
	v_sub_u32_e32 v34, vcc, 0, v26
	v_cmp_eq_u32_e64 s[2:3], 0, v32
	v_cndmask_b32_e64 v26, v26, v34, s[2:3]
	v_mul_hi_u32 v26, v26, v31
	v_add_u32_e32 v32, vcc, v33, v27
	v_xor_b32_e32 v32, v32, v33
	v_add_u32_e32 v34, vcc, v26, v31
	v_subrev_u32_e32 v26, vcc, v26, v31
	v_cndmask_b32_e64 v26, v26, v34, s[2:3]
	v_mul_hi_u32 v31, v26, v32
	v_add_u32_e32 v26, vcc, v28, v24
	v_mul_lo_u32 v28, v31, v30
	v_add_u32_e32 v33, vcc, -1, v31
	v_subrev_u32_e32 v34, vcc, v28, v32
	v_cmp_ge_u32_e64 s[4:5], v32, v28
	v_cmp_ge_u32_e64 s[2:3], v34, v30
	v_add_u32_e32 v28, vcc, 1, v31
	s_and_b64 vcc, s[2:3], s[4:5]
	v_cndmask_b32_e32 v28, v31, v28, vcc
	v_cndmask_b32_e64 v28, v33, v28, s[4:5]
	v_xor_b32_e32 v30, v28, v29
	v_sub_u32_e32 v28, vcc, v30, v29
	v_cmp_gt_i32_e32 vcc, 1, v28
	s_and_b64 vcc, exec, vcc
	s_cbranch_vccnz BB2_93
; %bb.89:                               ; %.preheader7.i.preheader
                                        ;   in Loop: Header=BB2_49 Depth=1
	v_not_b32_e32 v29, v29
	v_add_u32_e32 v29, vcc, v30, v29
	s_mov_b32 s2, 0
	v_mov_b32_e32 v30, v25
	s_branch BB2_91
BB2_90:                                 ;   in Loop: Header=BB2_91 Depth=2
	v_mov_b32_e32 v30, 0
	ds_read_b32 v30, v30 offset:36
	s_mov_b64 s[4:5], 0
	s_andn2_b64 vcc, exec, s[4:5]
	s_add_i32 s2, s2, 1
	s_cbranch_vccz BB2_93
BB2_91:                                 ; %.preheader7.i
                                        ;   Parent Loop BB2_49 Depth=1
                                        ; =>  This Inner Loop Header: Depth=2
	s_waitcnt lgkmcnt(0)
	v_mul_lo_u32 v30, s2, v30
	v_add_u32_e32 v30, vcc, v30, v0
	v_add_u32_e32 v31, vcc, v30, v24
	v_add_u32_e32 v30, vcc, v30, v26
	v_lshlrev_b32_e32 v31, 2, v31
	v_lshlrev_b32_e32 v32, 2, v30
	ds_read2_b32 v[30:31], v31 offset0:16 offset1:17
	ds_read_b32 v33, v32 offset:64
	s_waitcnt lgkmcnt(1)
	v_add_u32_e32 v30, vcc, v31, v30
	v_lshrrev_b32_e32 v31, 31, v30
	v_add_u32_e32 v30, vcc, v30, v31
	v_ashrrev_i32_e32 v30, 1, v30
	s_waitcnt lgkmcnt(0)
	v_sub_u32_e32 v30, vcc, v33, v30
	v_cmp_eq_u32_e32 vcc, s2, v29
	s_and_b64 vcc, exec, vcc
	ds_write_b32 v32, v30 offset:64
	s_cbranch_vccz BB2_90
; %bb.92:                               ;   in Loop: Header=BB2_91 Depth=2
	s_mov_b64 s[4:5], -1
                                        ; implicit-def: $vgpr30
	s_andn2_b64 vcc, exec, s[4:5]
	s_add_i32 s2, s2, 1
	s_cbranch_vccnz BB2_91
BB2_93:                                 ; %Flow385
                                        ;   in Loop: Header=BB2_49 Depth=1
	v_mul_lo_u32 v25, v28, v25
	v_subrev_u32_e32 v28, vcc, v25, v27
	v_subrev_u32_e32 v25, vcc, v28, v27
	v_ashrrev_i32_e32 v29, 31, v28
	v_cmp_lt_u64_e32 vcc, v[0:1], v[28:29]
	s_and_saveexec_b64 s[2:3], vcc
	s_cbranch_execz BB2_95
; %bb.94:                               ;   in Loop: Header=BB2_49 Depth=1
	v_add_u32_e32 v25, vcc, v25, v0
	v_add_u32_e32 v24, vcc, v25, v24
	v_add_u32_e32 v25, vcc, v25, v26
	v_lshlrev_b32_e32 v24, 2, v24
	v_lshlrev_b32_e32 v26, 2, v25
	ds_read2_b32 v[24:25], v24 offset0:16 offset1:17
	ds_read_b32 v27, v26 offset:64
	s_waitcnt lgkmcnt(1)
	v_add_u32_e32 v24, vcc, v25, v24
	v_lshrrev_b32_e32 v25, 31, v24
	v_add_u32_e32 v24, vcc, v24, v25
	v_ashrrev_i32_e32 v24, 1, v24
	s_waitcnt lgkmcnt(0)
	v_sub_u32_e32 v24, vcc, v27, v24
	ds_write_b32 v26, v24 offset:64
BB2_95:                                 ;   in Loop: Header=BB2_49 Depth=1
	s_or_b64 exec, exec, s[2:3]
	v_mov_b32_e32 v26, 0
	s_waitcnt lgkmcnt(0)
	s_barrier
	s_waitcnt lgkmcnt(0)
	ds_read2_b32 v[24:25], v26 offset0:1 offset1:9
	s_waitcnt lgkmcnt(0)
	v_ashrrev_i32_e32 v29, 31, v25
	v_add_u32_e32 v27, vcc, v29, v25
	v_xor_b32_e32 v30, v27, v29
	v_cvt_f32_u32_e32 v27, v30
	v_rcp_iflag_f32_e32 v31, v27
	ds_read2_b32 v[27:28], v26 offset0:11 offset1:15
	v_mul_f32_e32 v26, s25, v31
	v_cvt_u32_f32_e32 v31, v26
	s_waitcnt lgkmcnt(0)
	v_mul_lo_u32 v24, v27, v24
	v_lshlrev_b32_e32 v27, 1, v27
	v_mul_hi_u32 v32, v31, v30
	v_add_u32_e32 v26, vcc, -1, v24
	v_mul_lo_u32 v24, v31, v30
	v_ashrrev_i32_e32 v33, 31, v26
	v_cmp_eq_u32_e64 s[2:3], 0, v32
	v_add_u32_e32 v32, vcc, v33, v26
	v_sub_u32_e32 v34, vcc, 0, v24
	v_cndmask_b32_e64 v24, v24, v34, s[2:3]
	v_mul_hi_u32 v24, v24, v31
	v_xor_b32_e32 v32, v32, v33
	v_xor_b32_e32 v29, v33, v29
	v_add_u32_e32 v34, vcc, v24, v31
	v_subrev_u32_e32 v24, vcc, v24, v31
	v_cndmask_b32_e64 v24, v24, v34, s[2:3]
	v_mul_hi_u32 v31, v24, v32
	v_add_u32_e32 v24, vcc, v28, v27
	v_mul_lo_u32 v28, v31, v30
	v_add_u32_e32 v33, vcc, -1, v31
	v_subrev_u32_e32 v34, vcc, v28, v32
	v_cmp_ge_u32_e64 s[4:5], v32, v28
	v_cmp_ge_u32_e64 s[2:3], v34, v30
	v_add_u32_e32 v28, vcc, 1, v31
	s_and_b64 vcc, s[2:3], s[4:5]
	v_cndmask_b32_e32 v28, v31, v28, vcc
	v_cndmask_b32_e64 v28, v33, v28, s[4:5]
	v_xor_b32_e32 v30, v28, v29
	v_sub_u32_e32 v28, vcc, v30, v29
	v_cmp_gt_i32_e32 vcc, 1, v28
	s_and_b64 vcc, exec, vcc
	s_cbranch_vccnz BB2_100
; %bb.96:                               ; %.preheader5.i.preheader
                                        ;   in Loop: Header=BB2_49 Depth=1
	v_not_b32_e32 v29, v29
	v_add_u32_e32 v29, vcc, v30, v29
	s_mov_b32 s2, 0
	v_mov_b32_e32 v30, v25
	s_branch BB2_98
BB2_97:                                 ;   in Loop: Header=BB2_98 Depth=2
	v_mov_b32_e32 v30, 0
	ds_read_b32 v30, v30 offset:36
	s_mov_b64 s[4:5], 0
	s_andn2_b64 vcc, exec, s[4:5]
	s_add_i32 s2, s2, 1
	s_cbranch_vccz BB2_100
BB2_98:                                 ; %.preheader5.i
                                        ;   Parent Loop BB2_49 Depth=1
                                        ; =>  This Inner Loop Header: Depth=2
	s_waitcnt lgkmcnt(0)
	v_mul_lo_u32 v30, s2, v30
	v_add_u32_e32 v30, vcc, v30, v0
	v_add_u32_e32 v31, vcc, v30, v24
	v_add_u32_e32 v30, vcc, v30, v27
	v_lshlrev_b32_e32 v31, 2, v31
	v_lshlrev_b32_e32 v32, 2, v30
	ds_read2_b32 v[30:31], v31 offset0:16 offset1:17
	ds_read_b32 v33, v32 offset:68
	s_waitcnt lgkmcnt(1)
	v_add_u32_e32 v30, vcc, v30, v31
	v_add_u32_e32 v30, vcc, 2, v30
	v_ashrrev_i32_e32 v31, 31, v30
	v_lshrrev_b32_e32 v31, 30, v31
	v_add_u32_e32 v30, vcc, v30, v31
	v_ashrrev_i32_e32 v30, 2, v30
	s_waitcnt lgkmcnt(0)
	v_add_u32_e32 v30, vcc, v30, v33
	v_cmp_eq_u32_e32 vcc, s2, v29
	s_and_b64 vcc, exec, vcc
	ds_write_b32 v32, v30 offset:68
	s_cbranch_vccz BB2_97
; %bb.99:                               ;   in Loop: Header=BB2_98 Depth=2
	s_mov_b64 s[4:5], -1
                                        ; implicit-def: $vgpr30
	s_andn2_b64 vcc, exec, s[4:5]
	s_add_i32 s2, s2, 1
	s_cbranch_vccnz BB2_98
BB2_100:                                ; %Flow382
                                        ;   in Loop: Header=BB2_49 Depth=1
	v_mul_lo_u32 v28, v28, v25
	v_or_b32_e32 v25, 1, v27
	v_subrev_u32_e32 v27, vcc, v28, v26
	v_subrev_u32_e32 v26, vcc, v27, v26
	v_ashrrev_i32_e32 v28, 31, v27
	v_cmp_lt_u64_e32 vcc, v[0:1], v[27:28]
	s_and_saveexec_b64 s[2:3], vcc
	s_cbranch_execz BB2_102
; %bb.101:                              ;   in Loop: Header=BB2_49 Depth=1
	v_add_u32_e32 v26, vcc, v26, v0
	v_add_u32_e32 v24, vcc, v26, v24
	v_add_u32_e32 v25, vcc, v26, v25
	v_lshlrev_b32_e32 v24, 2, v24
	v_lshlrev_b32_e32 v26, 2, v25
	ds_read2_b32 v[24:25], v24 offset0:16 offset1:17
	ds_read_b32 v27, v26 offset:64
	s_waitcnt lgkmcnt(1)
	v_add_u32_e32 v24, vcc, v24, v25
	v_add_u32_e32 v24, vcc, 2, v24
	v_ashrrev_i32_e32 v25, 31, v24
	v_lshrrev_b32_e32 v25, 30, v25
	v_add_u32_e32 v24, vcc, v24, v25
	v_ashrrev_i32_e32 v24, 2, v24
	s_waitcnt lgkmcnt(0)
	v_add_u32_e32 v24, vcc, v24, v27
	ds_write_b32 v26, v24 offset:64
BB2_102:                                ;   in Loop: Header=BB2_49 Depth=1
	s_or_b64 exec, exec, s[2:3]
	v_mov_b32_e32 v24, 0
	s_waitcnt lgkmcnt(0)
	s_barrier
	s_waitcnt lgkmcnt(0)
	ds_read_b32 v25, v24 offset:4
	s_waitcnt lgkmcnt(0)
	v_cmp_gt_i32_e32 vcc, 1, v25
	s_and_b64 vcc, exec, vcc
	s_cbranch_vccnz BB2_48
; %bb.103:                              ;   in Loop: Header=BB2_49 Depth=1
	v_add_u32_e32 v25, vcc, -1, v25
	v_lshrrev_b32_e32 v26, 1, v25
	v_add_u32_e32 v26, vcc, 1, v26
	v_cmp_gt_u32_e32 vcc, 6, v25
	ds_read_b32 v25, v24 offset:44
	s_and_b64 vcc, exec, vcc
	s_cbranch_vccnz BB2_122
; %bb.104:                              ;   in Loop: Header=BB2_49 Depth=1
	s_waitcnt lgkmcnt(0)
	v_mul_lo_u32 v24, v25, 9
	v_mul_lo_u32 v30, v25, 7
	v_mul_lo_u32 v31, v25, 6
	v_mul_lo_u32 v34, v25, 12
	v_add_u32_e32 v24, vcc, v21, v24
	v_lshlrev_b32_e32 v28, 2, v24
	v_lshlrev_b32_e32 v24, 3, v25
	v_add_u32_e32 v24, vcc, v21, v24
	v_lshlrev_b32_e32 v29, 2, v24
	v_add_u32_e32 v24, vcc, v21, v30
	v_lshlrev_b32_e32 v30, 2, v24
	v_mul_lo_u32 v24, v25, 5
	v_add_u32_e32 v31, vcc, v21, v31
	s_mov_b32 s27, 2
	v_lshlrev_b32_e32 v27, 5, v25
	v_add_u32_e32 v24, vcc, v21, v24
	v_lshlrev_b32_e32 v32, 2, v24
	v_lshlrev_b32_e32 v24, 2, v25
	v_add_u32_e32 v24, vcc, v21, v24
	v_lshlrev_b32_e32 v33, 2, v24
	v_lshlrev_b32_e32 v24, 1, v25
	v_add_u32_e32 v24, vcc, v21, v24
	v_lshlrev_b32_e32 v35, 2, v24
	v_and_b32_e32 v24, s24, v26
	v_lshlrev_b32_e32 v31, 2, v31
	v_add_u32_e32 v34, vcc, v22, v34
	v_sub_u32_e32 v36, vcc, 0, v24
	v_mov_b32_e32 v37, 0
	s_branch BB2_107
BB2_105:                                ; %Flow373
                                        ;   in Loop: Header=BB2_107 Depth=2
	s_or_b64 exec, exec, s[4:5]
BB2_106:                                ;   in Loop: Header=BB2_107 Depth=2
	s_or_b64 exec, exec, s[28:29]
	v_cndmask_b32_e64 v15, 0, 1, vcc
	v_add_u32_e32 v37, vcc, v37, v27
	v_add_u32_e32 v36, vcc, 4, v36
	v_cmp_eq_u32_e32 vcc, 0, v36
	s_and_b64 vcc, exec, vcc
	s_add_i32 s27, s27, 8
	s_cbranch_vccnz BB2_123
BB2_107:                                ;   Parent Loop BB2_49 Depth=1
                                        ; =>  This Inner Loop Header: Depth=2
	v_and_b32_e32 v15, s26, v15
	v_cmp_eq_u16_e64 s[2:3], 0, v15
	v_cmp_ne_u32_e64 s[4:5], v23, v16
	v_cmp_ne_u16_e32 vcc, 0, v15
	s_or_b64 s[2:3], s[2:3], s[4:5]
	v_mov_b32_e32 v38, v16
	s_and_saveexec_b64 s[4:5], s[2:3]
	s_cbranch_execz BB2_111
; %bb.108:                              ;   in Loop: Header=BB2_107 Depth=2
	v_add_u32_e64 v15, s[2:3], v37, v35
	v_ashrrev_i32_e32 v24, 31, v23
	ds_read_b32 v15, v15 offset:64
	v_lshlrev_b64 v[39:40], 2, v[23:24]
	v_mov_b32_e32 v24, s11
	v_add_u32_e64 v39, s[2:3], s10, v39
	v_addc_u32_e64 v40, s[2:3], v24, v40, s[2:3]
	s_waitcnt lgkmcnt(0)
	flat_store_dword v[39:40], v15
	v_add_u32_e64 v15, s[2:3], v23, v14
	v_cmp_ne_u32_e64 s[2:3], v15, v16
	s_xor_b64 s[28:29], vcc, -1
	s_or_b64 s[2:3], s[28:29], s[2:3]
	v_mov_b32_e32 v38, v16
	s_and_saveexec_b64 s[28:29], s[2:3]
	s_cbranch_execz BB2_110
; %bb.109:                              ;   in Loop: Header=BB2_107 Depth=2
	v_add_u32_e64 v23, s[2:3], v37, v34
	ds_read_b32 v38, v23 offset:64
	v_ashrrev_i32_e32 v24, 31, v14
	v_mov_b32_e32 v23, v14
	v_lshlrev_b64 v[23:24], 2, v[23:24]
	v_add_u32_e64 v23, s[2:3], v39, v23
	v_addc_u32_e64 v24, s[2:3], v40, v24, s[2:3]
	s_waitcnt lgkmcnt(0)
	flat_store_dword v[23:24], v38
	v_add_u32_e64 v38, s[2:3], v15, v17
BB2_110:                                ; %Flow376
                                        ;   in Loop: Header=BB2_107 Depth=2
	s_or_b64 exec, exec, s[28:29]
BB2_111:                                ;   in Loop: Header=BB2_107 Depth=2
	s_or_b64 exec, exec, s[4:5]
	v_cmp_ne_u32_e64 s[2:3], v38, v16
	s_xor_b64 s[4:5], vcc, -1
	s_or_b64 s[2:3], s[4:5], s[2:3]
	v_mov_b32_e32 v23, v16
	s_and_saveexec_b64 s[28:29], s[2:3]
	s_cbranch_execz BB2_115
; %bb.112:                              ;   in Loop: Header=BB2_107 Depth=2
	v_add_u32_e64 v15, s[2:3], v37, v33
	v_ashrrev_i32_e32 v39, 31, v38
	ds_read_b32 v15, v15 offset:64
	v_lshlrev_b64 v[23:24], 2, v[38:39]
	v_mov_b32_e32 v40, s11
	v_add_u32_e64 v39, s[2:3], s10, v23
	v_addc_u32_e64 v40, s[2:3], v40, v24, s[2:3]
	s_waitcnt lgkmcnt(0)
	flat_store_dword v[39:40], v15
	v_add_u32_e64 v15, s[2:3], v38, v14
	v_cmp_ne_u32_e64 s[2:3], v15, v16
	s_or_b64 s[2:3], s[4:5], s[2:3]
	v_mov_b32_e32 v23, v16
	s_and_saveexec_b64 s[30:31], s[2:3]
	s_cbranch_execz BB2_114
; %bb.113:                              ;   in Loop: Header=BB2_107 Depth=2
	v_add_u32_e64 v23, s[2:3], v37, v32
	ds_read_b32 v38, v23 offset:64
	v_ashrrev_i32_e32 v24, 31, v14
	v_mov_b32_e32 v23, v14
	v_lshlrev_b64 v[23:24], 2, v[23:24]
	v_add_u32_e64 v23, s[2:3], v39, v23
	v_addc_u32_e64 v24, s[2:3], v40, v24, s[2:3]
	s_waitcnt lgkmcnt(0)
	flat_store_dword v[23:24], v38
	v_add_u32_e64 v23, s[2:3], v15, v17
BB2_114:                                ; %Flow375
                                        ;   in Loop: Header=BB2_107 Depth=2
	s_or_b64 exec, exec, s[30:31]
BB2_115:                                ;   in Loop: Header=BB2_107 Depth=2
	s_or_b64 exec, exec, s[28:29]
	v_cmp_ne_u32_e64 s[2:3], v23, v16
	s_or_b64 s[2:3], s[4:5], s[2:3]
	v_mov_b32_e32 v38, v16
	s_and_saveexec_b64 s[28:29], s[2:3]
	s_cbranch_execz BB2_119
; %bb.116:                              ;   in Loop: Header=BB2_107 Depth=2
	v_add_u32_e64 v15, s[2:3], v37, v31
	v_ashrrev_i32_e32 v24, 31, v23
	ds_read_b32 v15, v15 offset:64
	v_lshlrev_b64 v[39:40], 2, v[23:24]
	v_mov_b32_e32 v24, s11
	v_add_u32_e64 v39, s[2:3], s10, v39
	v_addc_u32_e64 v40, s[2:3], v24, v40, s[2:3]
	s_waitcnt lgkmcnt(0)
	flat_store_dword v[39:40], v15
	v_add_u32_e64 v15, s[2:3], v23, v14
	v_cmp_ne_u32_e64 s[2:3], v15, v16
	s_or_b64 s[2:3], s[4:5], s[2:3]
	v_mov_b32_e32 v38, v16
	s_and_saveexec_b64 s[30:31], s[2:3]
	s_cbranch_execz BB2_118
; %bb.117:                              ;   in Loop: Header=BB2_107 Depth=2
	v_add_u32_e64 v23, s[2:3], v37, v30
	ds_read_b32 v38, v23 offset:64
	v_ashrrev_i32_e32 v24, 31, v14
	v_mov_b32_e32 v23, v14
	v_lshlrev_b64 v[23:24], 2, v[23:24]
	v_add_u32_e64 v23, s[2:3], v39, v23
	v_addc_u32_e64 v24, s[2:3], v40, v24, s[2:3]
	s_waitcnt lgkmcnt(0)
	flat_store_dword v[23:24], v38
	v_add_u32_e64 v38, s[2:3], v15, v17
BB2_118:                                ; %Flow374
                                        ;   in Loop: Header=BB2_107 Depth=2
	s_or_b64 exec, exec, s[30:31]
BB2_119:                                ;   in Loop: Header=BB2_107 Depth=2
	s_or_b64 exec, exec, s[28:29]
	v_cmp_ne_u32_e64 s[2:3], v38, v16
	s_or_b64 s[2:3], s[4:5], s[2:3]
	v_mov_b32_e32 v23, v16
	s_and_saveexec_b64 s[28:29], s[2:3]
	s_cbranch_execz BB2_106
; %bb.120:                              ;   in Loop: Header=BB2_107 Depth=2
	v_add_u32_e64 v15, s[2:3], v37, v29
	v_ashrrev_i32_e32 v39, 31, v38
	ds_read_b32 v15, v15 offset:64
	v_lshlrev_b64 v[23:24], 2, v[38:39]
	v_mov_b32_e32 v40, s11
	v_add_u32_e64 v39, s[2:3], s10, v23
	v_addc_u32_e64 v40, s[2:3], v40, v24, s[2:3]
	s_waitcnt lgkmcnt(0)
	flat_store_dword v[39:40], v15
	v_add_u32_e64 v15, s[2:3], v38, v14
	v_cmp_ne_u32_e64 s[2:3], v15, v16
	s_or_b64 s[2:3], s[4:5], s[2:3]
	v_mov_b32_e32 v23, v16
	s_and_saveexec_b64 s[4:5], s[2:3]
	s_cbranch_execz BB2_105
; %bb.121:                              ;   in Loop: Header=BB2_107 Depth=2
	v_add_u32_e64 v23, s[2:3], v37, v28
	ds_read_b32 v38, v23 offset:64
	v_ashrrev_i32_e32 v24, 31, v14
	v_mov_b32_e32 v23, v14
	v_lshlrev_b64 v[23:24], 2, v[23:24]
	v_add_u32_e64 v23, s[2:3], v39, v23
	v_addc_u32_e64 v24, s[2:3], v40, v24, s[2:3]
	s_waitcnt lgkmcnt(0)
	flat_store_dword v[23:24], v38
	v_add_u32_e64 v23, s[2:3], v15, v17
	s_branch BB2_105
BB2_122:                                ;   in Loop: Header=BB2_49 Depth=1
	s_mov_b32 s27, 2
BB2_123:                                ; %Flow378
                                        ;   in Loop: Header=BB2_49 Depth=1
	v_and_b32_e32 v26, 3, v26
	v_cmp_eq_u32_e32 vcc, 0, v26
	s_and_b64 vcc, exec, vcc
	s_cbranch_vccnz BB2_48
; %bb.124:                              ; %.preheader.i.preheader
                                        ;   in Loop: Header=BB2_49 Depth=1
	s_waitcnt lgkmcnt(0)
	v_mul_lo_u32 v24, s27, v25
	v_lshlrev_b32_e32 v28, 3, v25
	v_add_u32_e32 v24, vcc, v21, v24
	v_lshlrev_b32_e32 v27, 2, v24
	s_branch BB2_127
BB2_125:                                ; %Flow
                                        ;   in Loop: Header=BB2_127 Depth=2
	s_or_b64 exec, exec, s[28:29]
BB2_126:                                ;   in Loop: Header=BB2_127 Depth=2
	s_or_b64 exec, exec, s[4:5]
	v_cndmask_b32_e64 v15, 0, 1, vcc
	v_add_u32_e32 v26, vcc, -1, v26
	v_add_u32_e32 v27, vcc, v27, v28
	v_cmp_ne_u32_e32 vcc, 0, v26
	s_add_i32 s27, s27, 2
	s_and_b64 vcc, exec, vcc
	v_mov_b32_e32 v23, v24
	s_cbranch_vccz BB2_47
BB2_127:                                ; %.preheader.i
                                        ;   Parent Loop BB2_49 Depth=1
                                        ; =>  This Inner Loop Header: Depth=2
	v_and_b32_e32 v15, s26, v15
	v_cmp_eq_u16_e64 s[2:3], 0, v15
	v_cmp_ne_u32_e64 s[4:5], v23, v16
	v_cmp_ne_u16_e32 vcc, 0, v15
	s_or_b64 s[2:3], s[2:3], s[4:5]
	v_mov_b32_e32 v24, v16
	s_and_saveexec_b64 s[4:5], s[2:3]
	s_cbranch_execz BB2_126
; %bb.128:                              ;   in Loop: Header=BB2_127 Depth=2
	v_ashrrev_i32_e32 v24, 31, v23
	v_lshlrev_b64 v[29:30], 2, v[23:24]
	ds_read_b32 v15, v27 offset:64
	v_mov_b32_e32 v24, s11
	v_add_u32_e64 v29, s[2:3], s10, v29
	v_addc_u32_e64 v30, s[2:3], v24, v30, s[2:3]
	v_add_u32_e64 v23, s[2:3], v23, v14
	v_cmp_ne_u32_e64 s[2:3], v23, v16
	s_xor_b64 s[28:29], vcc, -1
	s_or_b64 s[2:3], s[28:29], s[2:3]
	v_mov_b32_e32 v24, v16
	s_waitcnt lgkmcnt(0)
	flat_store_dword v[29:30], v15
	s_and_saveexec_b64 s[28:29], s[2:3]
	s_cbranch_execz BB2_125
; %bb.129:                              ;   in Loop: Header=BB2_127 Depth=2
	s_or_b32 s2, s27, 1
	v_mul_lo_u32 v24, s2, v25
	v_ashrrev_i32_e32 v15, 31, v14
	v_lshlrev_b64 v[31:32], 2, v[14:15]
	v_add_u32_e64 v15, s[2:3], v24, v18
	v_lshlrev_b32_e32 v15, 2, v15
	ds_read_b32 v15, v15 offset:64
	v_add_u32_e64 v29, s[2:3], v29, v31
	v_addc_u32_e64 v30, s[2:3], v30, v32, s[2:3]
	v_add_u32_e64 v24, s[2:3], v23, v17
	s_waitcnt lgkmcnt(0)
	flat_store_dword v[29:30], v15
	s_branch BB2_125
BB2_130:                                ; %Flow414
	s_mov_b64 s[2:3], 0
BB2_131:                                ; %Flow543
	s_and_b64 vcc, exec, s[2:3]
	s_cbranch_vccz BB2_385
; %bb.132:
	s_add_i32 s0, s6, 1
	s_mul_i32 s0, s0, s15
	s_add_i32 s0, s0, 1
	s_cmp_lt_i32 s0, s12
	s_mov_b64 s[0:1], -1
	s_cbranch_scc1 BB2_259
; %bb.133:
	v_add_u32_e32 v2, vcc, s23, v0
	v_cmp_ne_u32_e64 s[2:3], s7, 0
	s_and_b64 vcc, exec, s[2:3]
	s_cbranch_vccz BB2_139
; %bb.134:
	v_cmp_le_i32_e32 vcc, s12, v2
                                        ; implicit-def: $vgpr1
	s_and_saveexec_b64 s[0:1], vcc
	s_xor_b64 s[0:1], exec, s[0:1]
; %bb.135:
	s_lshl_b32 s4, s12, 1
	v_sub_u32_e32 v1, vcc, s4, v2
	v_add_u32_e32 v1, vcc, -2, v1
; %bb.136:                              ; %Flow479
	s_or_saveexec_b64 s[0:1], s[0:1]
	s_xor_b64 exec, exec, s[0:1]
; %bb.137:
	v_ashrrev_i32_e32 v1, 31, v2
	v_add_u32_e32 v3, vcc, v2, v1
	v_xor_b32_e32 v1, v3, v1
; %bb.138:
	s_or_b64 exec, exec, s[0:1]
	s_add_i32 s0, s21, -2
	s_mul_i32 s0, s0, s12
	v_add_u32_e32 v3, vcc, s0, v1
	s_add_i32 s0, s12, s12
	v_ashrrev_i32_e32 v4, 31, v3
	v_add_u32_e32 v1, vcc, s0, v3
	v_lshlrev_b64 v[3:4], 2, v[3:4]
	s_ashr_i32 s17, s12, 31
	s_waitcnt vmcnt(0) lgkmcnt(0)
	v_mov_b32_e32 v5, s9
	v_add_u32_e32 v3, vcc, s8, v3
	v_addc_u32_e32 v4, vcc, v5, v4, vcc
	s_lshl_b64 s[0:1], s[16:17], 2
	v_mov_b32_e32 v5, s1
	v_add_u32_e32 v7, vcc, s0, v3
	v_addc_u32_e32 v8, vcc, v4, v5, vcc
	v_add_u32_e32 v9, vcc, s0, v7
	v_addc_u32_e32 v10, vcc, v8, v5, vcc
	flat_load_dword v6, v[3:4]
	flat_load_dword v5, v[7:8]
	flat_load_dword v4, v[9:10]
	s_mov_b64 s[0:1], 0
	s_and_b64 vcc, exec, s[0:1]
	s_cbranch_vccnz BB2_140
	s_branch BB2_149
BB2_139:
                                        ; implicit-def: $vgpr4
	s_waitcnt vmcnt(0) lgkmcnt(0)
                                        ; implicit-def: $vgpr5
                                        ; implicit-def: $vgpr1
                                        ; implicit-def: $vgpr6
	s_cbranch_execz BB2_149
BB2_140:
	v_cmp_le_i32_e32 vcc, s12, v2
                                        ; implicit-def: $vgpr1
	s_and_saveexec_b64 s[0:1], vcc
	s_xor_b64 s[4:5], exec, s[0:1]
; %bb.141:
	s_lshl_b32 s0, s12, 1
	v_sub_u32_e64 v1, s[0:1], s0, v2
	v_add_u32_e64 v1, s[0:1], -2, v1
; %bb.142:                              ; %Flow481
	s_or_saveexec_b64 s[4:5], s[4:5]
	s_xor_b64 exec, exec, s[4:5]
; %bb.143:
	v_ashrrev_i32_e32 v1, 31, v2
	v_add_u32_e64 v3, s[0:1], v2, v1
	v_xor_b32_e32 v1, v3, v1
; %bb.144:
	s_or_b64 exec, exec, s[4:5]
	s_mul_i32 s4, s21, s12
	v_add_u32_e64 v3, s[0:1], s4, v1
	s_waitcnt vmcnt(0) lgkmcnt(0)
	v_ashrrev_i32_e32 v4, 31, v3
	v_lshlrev_b64 v[3:4], 2, v[3:4]
	s_ashr_i32 s17, s12, 31
	v_mov_b32_e32 v1, s9
	v_add_u32_e64 v3, s[0:1], s8, v3
	v_addc_u32_e64 v4, s[0:1], v1, v4, s[0:1]
	s_lshl_b64 s[24:25], s[16:17], 2
	v_mov_b32_e32 v1, s25
	v_add_u32_e64 v5, s[0:1], s24, v3
	v_addc_u32_e64 v6, s[0:1], v4, v1, s[0:1]
	v_add_u32_e64 v7, s[0:1], s24, v5
	v_addc_u32_e64 v8, s[0:1], v6, v1, s[0:1]
	flat_load_dword v4, v[3:4]
	flat_load_dword v5, v[5:6]
	flat_load_dword v6, v[7:8]
                                        ; implicit-def: $vgpr1
	s_and_saveexec_b64 s[0:1], vcc
	s_xor_b64 s[0:1], exec, s[0:1]
; %bb.145:
	s_lshl_b32 s5, s12, 1
	v_sub_u32_e32 v1, vcc, s5, v2
	v_add_u32_e32 v1, vcc, -2, v1
; %bb.146:                              ; %Flow480
	s_or_saveexec_b64 s[0:1], s[0:1]
	s_xor_b64 exec, exec, s[0:1]
; %bb.147:
	v_ashrrev_i32_e32 v1, 31, v2
	v_add_u32_e32 v2, vcc, v2, v1
	v_xor_b32_e32 v1, v2, v1
; %bb.148:
	s_or_b64 exec, exec, s[0:1]
	v_add_u32_e32 v1, vcc, s4, v1
BB2_149:                                ; %Flow483
	v_add_u32_e32 v7, vcc, 2, v0
	v_and_b32_e32 v2, 1, v7
	v_mul_lo_u32 v8, s20, v2
	v_mov_b32_e32 v3, 0
	v_lshrrev_b32_e32 v7, 1, v7
	v_mov_b32_e32 v2, v0
	v_add_u32_e32 v9, vcc, v8, v7
	v_cmp_gt_u32_e64 s[0:1], 3, v2
	v_mov_b32_e32 v10, v3
	v_mov_b32_e32 v11, v3
	v_mov_b32_e32 v12, v3
	v_mov_b32_e32 v14, v3
	v_mov_b32_e32 v7, v3
	v_mov_b32_e32 v8, v3
	s_and_saveexec_b64 s[4:5], s[0:1]
	s_cbranch_execz BB2_167
; %bb.150:
	v_mov_b32_e32 v7, s15
	v_cmp_eq_u32_e32 vcc, 0, v2
	v_cndmask_b32_e32 v7, -3, v7, vcc
	v_add_u32_e32 v8, vcc, v7, v2
	v_add_u32_e32 v13, vcc, s23, v8
	s_and_b64 vcc, exec, s[2:3]
	s_cbranch_vccz BB2_156
; %bb.151:
	v_cmp_le_i32_e32 vcc, s12, v13
                                        ; implicit-def: $vgpr7
	s_and_saveexec_b64 s[2:3], vcc
	s_xor_b64 s[2:3], exec, s[2:3]
; %bb.152:
	s_lshl_b32 s17, s12, 1
	v_sub_u32_e32 v7, vcc, s17, v13
	v_add_u32_e32 v7, vcc, -2, v7
; %bb.153:                              ; %Flow474
	s_or_saveexec_b64 s[2:3], s[2:3]
	s_xor_b64 exec, exec, s[2:3]
; %bb.154:
	v_ashrrev_i32_e32 v7, 31, v13
	v_add_u32_e32 v10, vcc, v13, v7
	v_xor_b32_e32 v7, v10, v7
; %bb.155:
	s_or_b64 exec, exec, s[2:3]
	s_add_i32 s2, s21, -2
	s_mul_i32 s2, s2, s12
	v_add_u32_e32 v10, vcc, s2, v7
	s_add_i32 s2, s12, s12
	v_ashrrev_i32_e32 v11, 31, v10
	v_add_u32_e32 v7, vcc, s2, v10
	v_lshlrev_b64 v[10:11], 2, v[10:11]
	v_mov_b32_e32 v12, s9
	v_add_u32_e32 v10, vcc, s8, v10
	s_ashr_i32 s3, s12, 31
	s_mov_b32 s2, s16
	v_addc_u32_e32 v11, vcc, v12, v11, vcc
	s_lshl_b64 s[2:3], s[2:3], 2
	v_mov_b32_e32 v12, s3
	v_add_u32_e32 v14, vcc, s2, v10
	v_addc_u32_e32 v15, vcc, v11, v12, vcc
	v_add_u32_e32 v16, vcc, s2, v14
	v_addc_u32_e32 v17, vcc, v15, v12, vcc
	flat_load_dword v12, v[10:11]
	flat_load_dword v11, v[14:15]
	flat_load_dword v10, v[16:17]
	s_mov_b64 s[2:3], 0
	s_and_b64 vcc, exec, s[2:3]
	s_cbranch_vccnz BB2_157
	s_branch BB2_166
BB2_156:
                                        ; implicit-def: $vgpr7
                                        ; implicit-def: $vgpr12
                                        ; implicit-def: $vgpr11
                                        ; implicit-def: $vgpr10
	s_cbranch_execz BB2_166
BB2_157:
	v_cmp_le_i32_e32 vcc, s12, v13
                                        ; implicit-def: $vgpr7
	s_and_saveexec_b64 s[2:3], vcc
	s_xor_b64 s[24:25], exec, s[2:3]
; %bb.158:
	s_lshl_b32 s2, s12, 1
	v_sub_u32_e64 v7, s[2:3], s2, v13
	v_add_u32_e64 v7, s[2:3], -2, v7
; %bb.159:                              ; %Flow476
	s_or_saveexec_b64 s[24:25], s[24:25]
	s_xor_b64 exec, exec, s[24:25]
	s_cbranch_execz BB2_161
; %bb.160:
	v_ashrrev_i32_e32 v7, 31, v13
	s_waitcnt vmcnt(0) lgkmcnt(0)
	v_add_u32_e64 v10, s[2:3], v13, v7
	v_xor_b32_e32 v7, v10, v7
BB2_161:
	s_or_b64 exec, exec, s[24:25]
	s_mul_i32 s17, s21, s12
	s_waitcnt vmcnt(0) lgkmcnt(0)
	v_add_u32_e64 v10, s[2:3], s17, v7
	v_ashrrev_i32_e32 v11, 31, v10
	v_lshlrev_b64 v[10:11], 2, v[10:11]
	v_mov_b32_e32 v7, s9
	v_add_u32_e64 v10, s[2:3], s8, v10
	v_addc_u32_e64 v11, s[2:3], v7, v11, s[2:3]
	s_ashr_i32 s3, s12, 31
	s_mov_b32 s2, s16
	s_lshl_b64 s[24:25], s[2:3], 2
	v_mov_b32_e32 v7, s25
	v_add_u32_e64 v14, s[2:3], s24, v10
	v_addc_u32_e64 v15, s[2:3], v11, v7, s[2:3]
	v_add_u32_e64 v16, s[2:3], s24, v14
	v_addc_u32_e64 v17, s[2:3], v15, v7, s[2:3]
	flat_load_dword v10, v[10:11]
	flat_load_dword v11, v[14:15]
	flat_load_dword v12, v[16:17]
                                        ; implicit-def: $vgpr7
	s_and_saveexec_b64 s[2:3], vcc
	s_xor_b64 s[2:3], exec, s[2:3]
; %bb.162:
	s_lshl_b32 s24, s12, 1
	v_sub_u32_e32 v7, vcc, s24, v13
	v_add_u32_e32 v7, vcc, -2, v7
; %bb.163:                              ; %Flow475
	s_or_saveexec_b64 s[2:3], s[2:3]
	s_xor_b64 exec, exec, s[2:3]
; %bb.164:
	v_ashrrev_i32_e32 v7, 31, v13
	v_add_u32_e32 v13, vcc, v13, v7
	v_xor_b32_e32 v7, v13, v7
; %bb.165:
	s_or_b64 exec, exec, s[2:3]
	v_add_u32_e32 v7, vcc, s17, v7
BB2_166:                                ; %Flow478
	v_add_u32_e32 v8, vcc, 2, v8
	v_and_b32_e32 v14, 1, v8
	v_mul_lo_u32 v14, v14, s20
	v_lshrrev_b32_e32 v13, 31, v8
	v_add_u32_e32 v8, vcc, v8, v13
	v_ashrrev_i32_e32 v8, 1, v8
	v_add_u32_e32 v14, vcc, v14, v8
	v_mov_b32_e32 v8, s12
BB2_167:
	s_or_b64 exec, exec, s[4:5]
	v_cvt_f32_u32_e32 v13, s22
	v_lshlrev_b64 v[17:18], 1, v[2:3]
	v_cvt_f32_u32_e32 v16, v2
	s_add_i32 s2, s15, -1
	v_rcp_iflag_f32_e32 v18, v13
	s_mul_i32 s3, s15, s6
	v_mov_b32_e32 v15, 0
	v_mov_b32_e32 v24, v15
	v_mul_f32_e32 v18, v16, v18
	v_trunc_f32_e32 v18, v18
	v_cvt_u32_f32_e32 v19, v18
	v_mad_f32 v16, -v18, v13, v16
	v_cmp_ge_f32_e64 vcc, |v16|, v13
	v_mov_b32_e32 v16, v15
	v_addc_u32_e32 v13, vcc, 0, v19, vcc
	v_and_b32_e32 v13, 0x3fffffff, v13
	v_mul_lo_u32 v18, v13, s2
	v_mov_b32_e32 v13, v15
	v_sub_u32_e32 v17, vcc, v17, v18
	v_add_u32_e32 v18, vcc, s3, v17
	v_cmp_gt_i32_e32 vcc, s12, v18
	s_and_saveexec_b64 s[2:3], vcc
	s_cbranch_execz BB2_173
; %bb.168:
	s_lshr_b32 s4, s12, 31
	s_add_i32 s4, s12, s4
	s_ashr_i32 s25, s4, 1
	s_and_b32 s4, s12, 1
	v_lshrrev_b32_e32 v13, 31, v18
	s_add_i32 s17, s25, s4
	s_lshr_b32 s4, s13, 31
	s_add_i32 s4, s13, s4
	v_add_u32_e32 v13, vcc, v18, v13
	v_ashrrev_i32_e32 v15, 1, v13
	v_and_b32_e32 v13, 1, v18
	s_ashr_i32 s4, s4, 1
	s_and_b32 s5, s13, 1
	v_cmp_eq_u32_e32 vcc, 1, v13
	s_add_i32 s24, s4, s5
                                        ; implicit-def: $vgpr18
                                        ; implicit-def: $sgpr28
	s_and_saveexec_b64 s[26:27], vcc
	s_xor_b64 s[26:27], exec, s[26:27]
; %bb.169:
	s_mul_i32 s28, s24, s17
	v_add_u32_e32 v18, vcc, s28, v15
	s_mul_i32 s28, s13, s12
	s_lshr_b32 s29, s28, 31
	s_add_i32 s28, s28, s29
	s_ashr_i32 s28, s28, 1
; %bb.170:                              ; %Flow472
	s_or_saveexec_b64 s[26:27], s[26:27]
	v_mov_b32_e32 v13, s28
	v_mov_b32_e32 v19, s25
	s_xor_b64 exec, exec, s[26:27]
; %bb.171:
	s_mul_i32 s24, s24, s12
	v_mov_b32_e32 v13, s24
	v_mov_b32_e32 v19, s17
	v_mov_b32_e32 v18, v15
; %bb.172:
	s_or_b64 exec, exec, s[26:27]
	v_mul_lo_u32 v15, v13, s5
	v_mul_lo_u32 v20, v19, s4
	s_lshr_b32 s4, s21, 31
	s_add_i32 s4, s21, s4
	v_add_u32_e32 v15, vcc, v15, v18
	s_and_b32 s5, s21, 1
	v_add_u32_e32 v15, vcc, v15, v20
	v_mul_lo_u32 v20, v13, s5
	s_ashr_i32 s4, s4, 1
	v_sub_u32_e32 v16, vcc, v19, v13
	v_mul_lo_u32 v19, v19, s4
	v_add_u32_e32 v18, vcc, v20, v18
	v_add_u32_e32 v24, vcc, v18, v19
BB2_173:                                ; %Flow473
	s_or_b64 exec, exec, s[2:3]
	v_add_u32_e32 v17, vcc, 2, v17
	v_and_b32_e32 v20, 1, v17
	v_mul_lo_u32 v19, v20, s20
	v_lshrrev_b32_e32 v18, 31, v17
	v_add_u32_e32 v17, vcc, v17, v18
	v_ashrrev_i32_e32 v23, 1, v17
	v_add_u32_e32 v17, vcc, v19, v23
	s_cmp_lt_i32 s14, 1
	s_cbranch_scc1 BB2_258
; %bb.174:
	s_sub_i32 s2, s18, s19
	v_lshlrev_b32_e32 v9, 2, v9
	s_add_i32 s2, s2, 32
	v_add_u32_e32 v18, vcc, 64, v9
	v_lshlrev_b32_e32 v9, 2, v14
	v_mul_lo_u32 v14, s2, v20
	v_add_u32_e32 v19, vcc, 64, v9
	v_ashrrev_i32_e32 v9, 31, v8
	s_ashr_i32 s17, s12, 31
	v_add_u32_e32 v22, vcc, v23, v14
	v_lshlrev_b64 v[20:21], 2, v[8:9]
	v_lshlrev_b32_e32 v14, 2, v14
	v_lshlrev_b32_e32 v23, 2, v23
	v_add_u32_e32 v23, vcc, v14, v23
	s_lshl_b64 s[24:25], s[16:17], 2
	s_mov_b32 s26, 0
	v_mov_b32_e32 v14, 1
	s_mov_b32 s27, 0x7ffffffc
	s_mov_b32 s28, 0x4f800000
	s_movk_i32 s29, 0xff
	s_branch BB2_177
BB2_175:                                ; %Flow428
                                        ;   in Loop: Header=BB2_177 Depth=1
	v_mov_b32_e32 v24, v25
BB2_176:                                ; %.loopexit.i49
                                        ;   in Loop: Header=BB2_177 Depth=1
	s_add_i32 s26, s26, 1
	s_cmp_eq_u32 s26, s14
	s_waitcnt vmcnt(0) lgkmcnt(0)
	s_barrier
	s_waitcnt lgkmcnt(0)
	s_cbranch_scc1 BB2_258
BB2_177:                                ; =>This Loop Header: Depth=1
                                        ;     Child Loop BB2_181 Depth 2
                                        ;     Child Loop BB2_185 Depth 2
                                        ;     Child Loop BB2_188 Depth 2
                                        ;     Child Loop BB2_192 Depth 2
                                        ;     Child Loop BB2_195 Depth 2
                                        ;     Child Loop BB2_201 Depth 2
                                        ;     Child Loop BB2_205 Depth 2
                                        ;     Child Loop BB2_208 Depth 2
                                        ;     Child Loop BB2_212 Depth 2
                                        ;     Child Loop BB2_215 Depth 2
                                        ;     Child Loop BB2_219 Depth 2
                                        ;     Child Loop BB2_226 Depth 2
                                        ;     Child Loop BB2_235 Depth 2
                                        ;     Child Loop BB2_255 Depth 2
	s_waitcnt vmcnt(0) lgkmcnt(0)
	ds_write_b32 v18, v6
	v_mov_b32_e32 v6, 0
	ds_read_b32 v25, v6 offset:8792
	s_mov_b64 s[2:3], -1
	s_waitcnt lgkmcnt(0)
	v_lshlrev_b32_e32 v25, 2, v25
	v_add_u32_e32 v25, vcc, v18, v25
	ds_write_b32 v25, v5
	ds_read_b32 v5, v6 offset:8792
	s_waitcnt lgkmcnt(0)
	v_lshlrev_b32_e32 v5, 3, v5
	v_add_u32_e32 v5, vcc, v18, v5
	ds_write_b32 v5, v4
	ds_read_b32 v4, v6 offset:4
                                        ; implicit-def: $vgpr5
	s_waitcnt lgkmcnt(0)
	v_cmp_gt_i32_e32 vcc, 1, v4
	s_and_b64 vcc, exec, vcc
	s_cbranch_vccz BB2_179
; %bb.178:                              ;   in Loop: Header=BB2_177 Depth=1
	v_add_u32_e32 v5, vcc, 2, v4
	s_mov_b64 s[2:3], 0
BB2_179:                                ; %Flow468
                                        ;   in Loop: Header=BB2_177 Depth=1
	s_andn2_b64 vcc, exec, s[2:3]
	s_cbranch_vccnz BB2_182
; %bb.180:                              ; %.preheader23.i31.preheader
                                        ;   in Loop: Header=BB2_177 Depth=1
	v_ashrrev_i32_e32 v5, 31, v1
	v_mov_b32_e32 v6, s17
	v_add_u32_e32 v4, vcc, s16, v1
	v_addc_u32_e32 v5, vcc, v6, v5, vcc
	v_lshlrev_b64 v[4:5], 2, v[4:5]
	v_mov_b32_e32 v6, s9
	v_add_u32_e32 v25, vcc, s8, v4
	s_mov_b32 s2, 2
	v_addc_u32_e32 v26, vcc, v6, v5, vcc
BB2_181:                                ; %.preheader23.i31
                                        ;   Parent Loop BB2_177 Depth=1
                                        ; =>  This Inner Loop Header: Depth=2
	flat_load_dword v4, v[25:26]
	v_mov_b32_e32 v5, 0
	ds_read_b32 v27, v5 offset:8792
	v_add_u32_e32 v1, vcc, s12, v1
	s_add_i32 s2, s2, 1
	v_mov_b32_e32 v6, s25
	v_add_u32_e32 v25, vcc, s24, v25
	v_addc_u32_e32 v26, vcc, v26, v6, vcc
	s_waitcnt lgkmcnt(0)
	v_mul_lo_u32 v6, v27, s2
	v_lshlrev_b32_e32 v6, 2, v6
	v_add_u32_e32 v6, vcc, v18, v6
	s_waitcnt vmcnt(0)
	ds_write_b32 v6, v4
	ds_read_b32 v4, v5 offset:4
	s_waitcnt lgkmcnt(0)
	v_add_u32_e32 v5, vcc, 2, v4
	v_cmp_ge_i32_e32 vcc, s2, v5
	s_and_b64 vcc, exec, vcc
	s_cbranch_vccz BB2_181
BB2_182:                                ; %.loopexit24.i27
                                        ;   in Loop: Header=BB2_177 Depth=1
	v_mov_b32_e32 v25, 0
	ds_read_b32 v6, v25 offset:8792
	ds_read_b32 v26, v25 offset:40
	s_waitcnt lgkmcnt(1)
	v_mul_lo_u32 v4, v6, v4
	v_mul_lo_u32 v5, v6, v5
	v_lshlrev_b32_e32 v6, 2, v6
	v_lshlrev_b32_e32 v4, 2, v4
	v_lshlrev_b32_e32 v5, 2, v5
	v_add_u32_e32 v27, vcc, v18, v4
	v_add_u32_e32 v4, vcc, v18, v5
	v_add_u32_e32 v5, vcc, v27, v6
	ds_read_b32 v4, v4
	ds_read_b32 v6, v27
	ds_read_b32 v5, v5
	s_waitcnt lgkmcnt(3)
	v_cmp_gt_i32_e32 vcc, 3, v26
	s_and_b64 vcc, exec, vcc
	s_cbranch_vccnz BB2_196
; %bb.183:                              ;   in Loop: Header=BB2_177 Depth=1
	v_add_u32_e32 v26, vcc, -1, v26
	v_lshrrev_b32_e32 v27, 31, v26
	v_add_u32_e32 v26, vcc, v27, v26
	v_ashrrev_i32_e32 v26, 1, v26
	v_max_i32_e32 v26, 1, v26
	v_add_u32_e32 v27, vcc, -1, v26
	v_cmp_gt_u32_e32 vcc, 3, v27
	s_and_b64 vcc, exec, vcc
	s_cbranch_vccnz BB2_186
; %bb.184:                              ;   in Loop: Header=BB2_177 Depth=1
	v_and_b32_e32 v25, s27, v26
	s_mov_b32 s4, 0
	s_mov_b32 s5, 8
BB2_185:                                ;   Parent Loop BB2_177 Depth=1
                                        ; =>  This Inner Loop Header: Depth=2
	v_mov_b32_e32 v27, 0
	ds_read_b32 v28, v27 offset:44
	s_add_i32 s30, s5, -8
	s_add_i32 s31, s5, -6
	s_add_i32 s33, s5, -4
	s_add_i32 s34, s5, -2
	s_waitcnt lgkmcnt(0)
	v_mul_lo_u32 v29, s31, v28
	v_mul_lo_u32 v30, v28, s30
	v_lshlrev_b32_e32 v28, 2, v28
	s_add_i32 s4, s4, 4
	v_lshlrev_b32_e32 v29, 2, v29
	v_lshlrev_b32_e32 v30, 2, v30
	v_add_u32_e64 v29, s[2:3], v18, v29
	v_add_u32_e64 v30, s[2:3], v18, v30
	ds_read_b32 v29, v29
	ds_read_b32 v31, v30
	v_add_u32_e64 v28, s[2:3], v30, v28
	ds_read_b32 v30, v28
	v_cmp_ne_u32_e32 vcc, s4, v25
	s_and_b64 vcc, exec, vcc
	s_waitcnt lgkmcnt(1)
	v_add_u32_e64 v29, s[2:3], v31, v29
	v_lshrrev_b32_e32 v31, 31, v29
	v_add_u32_e64 v29, s[2:3], v29, v31
	v_ashrrev_i32_e32 v29, 1, v29
	s_waitcnt lgkmcnt(0)
	v_sub_u32_e64 v29, s[2:3], v30, v29
	ds_write_b32 v28, v29
	ds_read_b32 v28, v27 offset:44
	s_waitcnt lgkmcnt(0)
	v_mul_lo_u32 v29, v28, s33
	v_mul_lo_u32 v30, v28, s31
	v_lshlrev_b32_e32 v28, 2, v28
	v_lshlrev_b32_e32 v29, 2, v29
	v_lshlrev_b32_e32 v30, 2, v30
	v_add_u32_e64 v29, s[2:3], v18, v29
	v_add_u32_e64 v30, s[2:3], v18, v30
	ds_read_b32 v29, v29
	ds_read_b32 v31, v30
	v_add_u32_e64 v28, s[2:3], v30, v28
	ds_read_b32 v30, v28
	s_waitcnt lgkmcnt(1)
	v_add_u32_e64 v29, s[2:3], v31, v29
	v_lshrrev_b32_e32 v31, 31, v29
	v_add_u32_e64 v29, s[2:3], v29, v31
	v_ashrrev_i32_e32 v29, 1, v29
	s_waitcnt lgkmcnt(0)
	v_sub_u32_e64 v29, s[2:3], v30, v29
	ds_write_b32 v28, v29
	ds_read_b32 v28, v27 offset:44
	s_waitcnt lgkmcnt(0)
	v_mul_lo_u32 v29, v28, s34
	v_mul_lo_u32 v30, v28, s33
	v_lshlrev_b32_e32 v28, 2, v28
	v_lshlrev_b32_e32 v29, 2, v29
	v_lshlrev_b32_e32 v30, 2, v30
	v_add_u32_e64 v29, s[2:3], v18, v29
	v_add_u32_e64 v30, s[2:3], v18, v30
	ds_read_b32 v29, v29
	ds_read_b32 v31, v30
	v_add_u32_e64 v28, s[2:3], v30, v28
	ds_read_b32 v30, v28
	s_waitcnt lgkmcnt(1)
	v_add_u32_e64 v29, s[2:3], v31, v29
	v_lshrrev_b32_e32 v31, 31, v29
	v_add_u32_e64 v29, s[2:3], v29, v31
	v_ashrrev_i32_e32 v29, 1, v29
	s_waitcnt lgkmcnt(0)
	v_sub_u32_e64 v29, s[2:3], v30, v29
	ds_write_b32 v28, v29
	ds_read_b32 v27, v27 offset:44
	s_waitcnt lgkmcnt(0)
	v_mul_lo_u32 v28, v27, s5
	v_mul_lo_u32 v29, v27, s34
	v_lshlrev_b32_e32 v27, 2, v27
	s_add_i32 s5, s5, 8
	v_lshlrev_b32_e32 v28, 2, v28
	v_lshlrev_b32_e32 v29, 2, v29
	v_add_u32_e64 v28, s[2:3], v18, v28
	v_add_u32_e64 v29, s[2:3], v18, v29
	ds_read_b32 v28, v28
	ds_read_b32 v30, v29
	v_add_u32_e64 v27, s[2:3], v29, v27
	ds_read_b32 v29, v27
	s_waitcnt lgkmcnt(1)
	v_add_u32_e64 v28, s[2:3], v30, v28
	v_lshrrev_b32_e32 v30, 31, v28
	v_add_u32_e64 v28, s[2:3], v28, v30
	v_ashrrev_i32_e32 v28, 1, v28
	s_waitcnt lgkmcnt(0)
	v_sub_u32_e64 v28, s[2:3], v29, v28
	ds_write_b32 v27, v28
	s_cbranch_vccnz BB2_185
BB2_186:                                ; %Flow465
                                        ;   in Loop: Header=BB2_177 Depth=1
	v_and_b32_e32 v26, 3, v26
	v_cmp_eq_u32_e32 vcc, 0, v26
	s_and_b64 vcc, exec, vcc
	s_cbranch_vccnz BB2_189
; %bb.187:                              ; %.preheader20.i29.preheader
                                        ;   in Loop: Header=BB2_177 Depth=1
	v_lshlrev_b32_e32 v25, 1, v25
BB2_188:                                ; %.preheader20.i29
                                        ;   Parent Loop BB2_177 Depth=1
                                        ; =>  This Inner Loop Header: Depth=2
	v_mov_b32_e32 v27, 0
	ds_read_b32 v27, v27 offset:44
	v_add_u32_e32 v28, vcc, 2, v25
	v_add_u32_e32 v26, vcc, -1, v26
	v_cmp_ne_u32_e32 vcc, 0, v26
	s_and_b64 vcc, exec, vcc
	s_waitcnt lgkmcnt(0)
	v_mul_lo_u32 v29, v28, v27
	v_mul_lo_u32 v30, v27, v25
	v_mov_b32_e32 v25, v28
	v_lshlrev_b32_e32 v27, 2, v27
	v_lshlrev_b32_e32 v29, 2, v29
	v_lshlrev_b32_e32 v28, 2, v30
	v_add_u32_e64 v29, s[2:3], v18, v29
	v_add_u32_e64 v28, s[2:3], v18, v28
	ds_read_b32 v29, v29
	ds_read_b32 v30, v28
	v_add_u32_e64 v27, s[2:3], v28, v27
	ds_read_b32 v28, v27
	s_waitcnt lgkmcnt(1)
	v_add_u32_e64 v29, s[2:3], v30, v29
	v_lshrrev_b32_e32 v30, 31, v29
	v_add_u32_e64 v29, s[2:3], v29, v30
	v_ashrrev_i32_e32 v29, 1, v29
	s_waitcnt lgkmcnt(0)
	v_sub_u32_e64 v28, s[2:3], v28, v29
	ds_write_b32 v27, v28
	s_cbranch_vccnz BB2_188
BB2_189:                                ; %.loopexit21.i30
                                        ;   in Loop: Header=BB2_177 Depth=1
	v_mov_b32_e32 v25, 0
	ds_read_b32 v26, v25 offset:40
	s_waitcnt lgkmcnt(0)
	v_cmp_gt_i32_e32 vcc, 4, v26
	s_and_b64 vcc, exec, vcc
	s_cbranch_vccnz BB2_196
; %bb.190:                              ;   in Loop: Header=BB2_177 Depth=1
	v_lshrrev_b32_e32 v26, 1, v26
	v_max_u32_e32 v27, 2, v26
	v_add_u32_e32 v26, vcc, -1, v27
	v_add_u32_e32 v27, vcc, -2, v27
	v_cmp_gt_u32_e32 vcc, 3, v27
	s_and_b64 vcc, exec, vcc
	s_cbranch_vccnz BB2_193
; %bb.191:                              ;   in Loop: Header=BB2_177 Depth=1
	v_and_b32_e32 v25, -4, v26
	s_mov_b32 s4, 0
	s_mov_b32 s5, 9
BB2_192:                                ;   Parent Loop BB2_177 Depth=1
                                        ; =>  This Inner Loop Header: Depth=2
	v_mov_b32_e32 v27, 0
	ds_read_b32 v28, v27 offset:44
	s_add_i32 s2, s5, -6
	s_add_i32 s3, s5, -4
	s_add_i32 s30, s5, -2
	s_add_i32 s4, s4, 4
	s_waitcnt lgkmcnt(0)
	v_mul_lo_u32 v29, v28, s2
	v_lshlrev_b32_e32 v31, 1, v28
	v_lshlrev_b32_e32 v28, 2, v28
	v_lshlrev_b32_e32 v30, 2, v29
	v_subrev_u32_e32 v29, vcc, v31, v29
	v_lshlrev_b32_e32 v29, 2, v29
	v_add_u32_e32 v30, vcc, v18, v30
	v_add_u32_e32 v29, vcc, v18, v29
	ds_read_b32 v30, v30
	ds_read_b32 v31, v29
	v_add_u32_e32 v28, vcc, v29, v28
	ds_read_b32 v29, v28
	s_waitcnt lgkmcnt(1)
	v_add_u32_e32 v30, vcc, v30, v31
	v_add_u32_e32 v30, vcc, 2, v30
	v_ashrrev_i32_e32 v31, 31, v30
	v_lshrrev_b32_e32 v31, 30, v31
	v_add_u32_e32 v30, vcc, v30, v31
	v_ashrrev_i32_e32 v30, 2, v30
	s_waitcnt lgkmcnt(0)
	v_add_u32_e32 v29, vcc, v30, v29
	ds_write_b32 v28, v29
	ds_read_b32 v28, v27 offset:44
	s_waitcnt lgkmcnt(0)
	v_mul_lo_u32 v31, v28, s2
	v_mul_lo_u32 v29, v28, s3
	v_lshlrev_b32_e32 v31, 2, v31
	v_lshlrev_b32_e32 v30, 2, v29
	v_add_u32_e32 v30, vcc, v18, v30
	v_add_u32_e32 v31, vcc, v18, v31
	ds_read_b32 v30, v30
	ds_read_b32 v31, v31
	v_subrev_u32_e32 v28, vcc, v28, v29
	v_lshlrev_b32_e32 v28, 2, v28
	v_add_u32_e32 v28, vcc, v18, v28
	ds_read_b32 v29, v28
	s_waitcnt lgkmcnt(1)
	v_add_u32_e32 v30, vcc, v30, v31
	v_add_u32_e32 v30, vcc, 2, v30
	v_ashrrev_i32_e32 v31, 31, v30
	v_lshrrev_b32_e32 v31, 30, v31
	v_add_u32_e32 v30, vcc, v30, v31
	v_ashrrev_i32_e32 v30, 2, v30
	s_waitcnt lgkmcnt(0)
	v_add_u32_e32 v29, vcc, v30, v29
	ds_write_b32 v28, v29
	ds_read_b32 v28, v27 offset:44
	s_waitcnt lgkmcnt(0)
	v_mul_lo_u32 v29, v28, s30
	v_lshlrev_b32_e32 v31, 1, v28
	v_lshlrev_b32_e32 v28, 2, v28
	v_lshlrev_b32_e32 v30, 2, v29
	v_subrev_u32_e32 v29, vcc, v31, v29
	v_lshlrev_b32_e32 v29, 2, v29
	v_add_u32_e32 v30, vcc, v18, v30
	v_add_u32_e32 v29, vcc, v18, v29
	ds_read_b32 v30, v30
	ds_read_b32 v31, v29
	v_add_u32_e32 v28, vcc, v29, v28
	ds_read_b32 v29, v28
	s_waitcnt lgkmcnt(1)
	v_add_u32_e32 v30, vcc, v30, v31
	v_add_u32_e32 v30, vcc, 2, v30
	v_ashrrev_i32_e32 v31, 31, v30
	v_lshrrev_b32_e32 v31, 30, v31
	v_add_u32_e32 v30, vcc, v30, v31
	v_ashrrev_i32_e32 v30, 2, v30
	s_waitcnt lgkmcnt(0)
	v_add_u32_e32 v29, vcc, v30, v29
	ds_write_b32 v28, v29
	ds_read_b32 v27, v27 offset:44
	s_waitcnt lgkmcnt(0)
	v_mul_lo_u32 v28, v27, s5
	v_mul_lo_u32 v30, v27, s30
	s_add_i32 s5, s5, 8
	v_lshlrev_b32_e32 v29, 2, v28
	v_lshlrev_b32_e32 v30, 2, v30
	v_add_u32_e32 v29, vcc, v18, v29
	v_add_u32_e32 v30, vcc, v18, v30
	ds_read_b32 v29, v29
	ds_read_b32 v30, v30
	v_subrev_u32_e32 v27, vcc, v27, v28
	v_lshlrev_b32_e32 v27, 2, v27
	v_add_u32_e64 v27, s[2:3], v18, v27
	ds_read_b32 v28, v27
	s_waitcnt lgkmcnt(1)
	v_add_u32_e64 v29, s[2:3], v29, v30
	v_add_u32_e64 v29, s[2:3], 2, v29
	v_ashrrev_i32_e32 v30, 31, v29
	v_lshrrev_b32_e32 v30, 30, v30
	v_add_u32_e64 v29, s[2:3], v29, v30
	v_cmp_eq_u32_e32 vcc, s4, v25
	v_ashrrev_i32_e32 v29, 2, v29
	s_and_b64 vcc, exec, vcc
	s_waitcnt lgkmcnt(0)
	v_add_u32_e64 v28, s[2:3], v29, v28
	ds_write_b32 v27, v28
	s_cbranch_vccz BB2_192
BB2_193:                                ; %Flow460
                                        ;   in Loop: Header=BB2_177 Depth=1
	v_and_b32_e32 v26, 3, v26
	v_cmp_eq_u32_e32 vcc, 0, v26
	s_and_b64 vcc, exec, vcc
	s_cbranch_vccnz BB2_196
; %bb.194:                              ; %.preheader17.i33.preheader
                                        ;   in Loop: Header=BB2_177 Depth=1
	v_lshlrev_b32_e32 v25, 1, v25
	v_add_u32_e32 v25, vcc, 3, v25
BB2_195:                                ; %.preheader17.i33
                                        ;   Parent Loop BB2_177 Depth=1
                                        ; =>  This Inner Loop Header: Depth=2
	v_mov_b32_e32 v27, 0
	ds_read_b32 v27, v27 offset:44
	v_add_u32_e32 v26, vcc, -1, v26
	s_waitcnt lgkmcnt(0)
	v_mul_lo_u32 v28, v27, v25
	v_lshlrev_b32_e32 v30, 1, v27
	v_lshlrev_b32_e32 v27, 2, v27
	v_add_u32_e32 v25, vcc, 2, v25
	v_lshlrev_b32_e32 v29, 2, v28
	v_subrev_u32_e32 v28, vcc, v30, v28
	v_lshlrev_b32_e32 v28, 2, v28
	v_add_u32_e32 v29, vcc, v18, v29
	v_add_u32_e64 v28, s[2:3], v18, v28
	ds_read_b32 v29, v29
	ds_read_b32 v30, v28
	v_add_u32_e64 v27, s[2:3], v28, v27
	ds_read_b32 v28, v27
	v_cmp_ne_u32_e32 vcc, 0, v26
	s_and_b64 vcc, exec, vcc
	s_waitcnt lgkmcnt(1)
	v_add_u32_e64 v29, s[2:3], v29, v30
	v_add_u32_e64 v29, s[2:3], 2, v29
	v_ashrrev_i32_e32 v30, 31, v29
	v_lshrrev_b32_e32 v30, 30, v30
	v_add_u32_e64 v29, s[2:3], v29, v30
	v_ashrrev_i32_e32 v29, 2, v29
	s_waitcnt lgkmcnt(0)
	v_add_u32_e64 v28, s[2:3], v29, v28
	ds_write_b32 v27, v28
	s_cbranch_vccnz BB2_195
BB2_196:                                ; %.loopexit18.i34
                                        ;   in Loop: Header=BB2_177 Depth=1
	s_and_saveexec_b64 s[4:5], s[0:1]
	s_cbranch_execz BB2_216
; %bb.197:                              ;   in Loop: Header=BB2_177 Depth=1
	ds_write_b32 v19, v12
	v_mov_b32_e32 v12, 0
	ds_read_b32 v25, v12 offset:8792
	s_waitcnt lgkmcnt(0)
	v_lshlrev_b32_e32 v25, 2, v25
	v_add_u32_e32 v25, vcc, v19, v25
	ds_write_b32 v25, v11
	ds_read_b32 v11, v12 offset:8792
	s_waitcnt lgkmcnt(0)
	v_lshlrev_b32_e32 v11, 3, v11
	v_add_u32_e32 v11, vcc, v19, v11
	ds_write_b32 v11, v10
	ds_read_b32 v10, v12 offset:4
	s_waitcnt lgkmcnt(0)
	v_cmp_gt_i32_e32 vcc, 1, v10
	s_and_b64 vcc, exec, vcc
	s_cbranch_vccz BB2_199
; %bb.198:                              ;   in Loop: Header=BB2_177 Depth=1
	v_add_u32_e32 v11, vcc, 2, v10
	s_mov_b64 s[2:3], 0
	s_andn2_b64 vcc, exec, s[2:3]
	s_cbranch_vccz BB2_200
	s_branch BB2_202
BB2_199:                                ;   in Loop: Header=BB2_177 Depth=1
	s_mov_b64 s[2:3], -1
                                        ; implicit-def: $vgpr11
	s_andn2_b64 vcc, exec, s[2:3]
	s_cbranch_vccnz BB2_202
BB2_200:                                ; %.preheader15.i39.preheader
                                        ;   in Loop: Header=BB2_177 Depth=1
	v_ashrrev_i32_e32 v11, 31, v7
	v_add_u32_e32 v10, vcc, v8, v7
	v_addc_u32_e32 v11, vcc, v9, v11, vcc
	v_lshlrev_b64 v[10:11], 2, v[10:11]
	v_mov_b32_e32 v12, s9
	v_add_u32_e32 v25, vcc, s8, v10
	s_mov_b32 s2, 2
	v_addc_u32_e32 v26, vcc, v12, v11, vcc
BB2_201:                                ; %.preheader15.i39
                                        ;   Parent Loop BB2_177 Depth=1
                                        ; =>  This Inner Loop Header: Depth=2
	flat_load_dword v10, v[25:26]
	v_mov_b32_e32 v11, 0
	ds_read_b32 v12, v11 offset:8792
	s_add_i32 s2, s2, 1
	v_add_u32_e32 v7, vcc, v7, v8
	v_add_u32_e32 v25, vcc, v25, v20
	v_addc_u32_e32 v26, vcc, v26, v21, vcc
	s_waitcnt lgkmcnt(0)
	v_mul_lo_u32 v12, v12, s2
	v_lshlrev_b32_e32 v12, 2, v12
	v_add_u32_e32 v12, vcc, v19, v12
	s_waitcnt vmcnt(0)
	ds_write_b32 v12, v10
	ds_read_b32 v10, v11 offset:4
	s_waitcnt lgkmcnt(0)
	v_add_u32_e32 v11, vcc, 2, v10
	v_cmp_ge_i32_e32 vcc, s2, v11
	s_and_b64 vcc, exec, vcc
	s_cbranch_vccz BB2_201
BB2_202:                                ; %.loopexit16.i35
                                        ;   in Loop: Header=BB2_177 Depth=1
	v_mov_b32_e32 v25, 0
	ds_read_b32 v12, v25 offset:8792
	ds_read_b32 v26, v25 offset:40
	s_waitcnt lgkmcnt(1)
	v_mul_lo_u32 v10, v12, v10
	v_mul_lo_u32 v11, v12, v11
	v_lshlrev_b32_e32 v12, 2, v12
	v_lshlrev_b32_e32 v10, 2, v10
	v_lshlrev_b32_e32 v11, 2, v11
	v_add_u32_e32 v27, vcc, v19, v10
	v_add_u32_e32 v10, vcc, v19, v11
	v_add_u32_e32 v11, vcc, v27, v12
	ds_read_b32 v10, v10
	ds_read_b32 v12, v27
	ds_read_b32 v11, v11
	s_waitcnt lgkmcnt(3)
	v_cmp_gt_i32_e32 vcc, 3, v26
	s_and_b64 vcc, exec, vcc
	s_cbranch_vccnz BB2_216
; %bb.203:                              ;   in Loop: Header=BB2_177 Depth=1
	v_add_u32_e32 v26, vcc, -1, v26
	v_lshrrev_b32_e32 v27, 31, v26
	v_add_u32_e32 v26, vcc, v27, v26
	v_ashrrev_i32_e32 v26, 1, v26
	v_max_i32_e32 v26, 1, v26
	v_add_u32_e32 v27, vcc, -1, v26
	v_cmp_gt_u32_e32 vcc, 3, v27
	s_and_b64 vcc, exec, vcc
	s_cbranch_vccnz BB2_206
; %bb.204:                              ;   in Loop: Header=BB2_177 Depth=1
	v_and_b32_e32 v25, s27, v26
	s_mov_b32 s30, 0
	s_mov_b32 s31, 8
BB2_205:                                ;   Parent Loop BB2_177 Depth=1
                                        ; =>  This Inner Loop Header: Depth=2
	v_mov_b32_e32 v27, 0
	ds_read_b32 v28, v27 offset:44
	s_add_i32 s33, s31, -8
	s_add_i32 s34, s31, -6
	s_add_i32 s35, s31, -4
	s_add_i32 s36, s31, -2
	s_waitcnt lgkmcnt(0)
	v_mul_lo_u32 v29, s34, v28
	v_mul_lo_u32 v30, v28, s33
	v_lshlrev_b32_e32 v28, 2, v28
	s_add_i32 s30, s30, 4
	v_lshlrev_b32_e32 v29, 2, v29
	v_lshlrev_b32_e32 v30, 2, v30
	v_add_u32_e64 v29, s[2:3], v19, v29
	v_add_u32_e64 v30, s[2:3], v19, v30
	ds_read_b32 v29, v29
	ds_read_b32 v31, v30
	v_add_u32_e64 v28, s[2:3], v30, v28
	ds_read_b32 v30, v28
	v_cmp_ne_u32_e32 vcc, s30, v25
	s_and_b64 vcc, exec, vcc
	s_waitcnt lgkmcnt(1)
	v_add_u32_e64 v29, s[2:3], v31, v29
	v_lshrrev_b32_e32 v31, 31, v29
	v_add_u32_e64 v29, s[2:3], v29, v31
	v_ashrrev_i32_e32 v29, 1, v29
	s_waitcnt lgkmcnt(0)
	v_sub_u32_e64 v29, s[2:3], v30, v29
	ds_write_b32 v28, v29
	ds_read_b32 v28, v27 offset:44
	s_waitcnt lgkmcnt(0)
	v_mul_lo_u32 v29, v28, s35
	v_mul_lo_u32 v30, v28, s34
	v_lshlrev_b32_e32 v28, 2, v28
	v_lshlrev_b32_e32 v29, 2, v29
	v_lshlrev_b32_e32 v30, 2, v30
	v_add_u32_e64 v29, s[2:3], v19, v29
	v_add_u32_e64 v30, s[2:3], v19, v30
	ds_read_b32 v29, v29
	ds_read_b32 v31, v30
	v_add_u32_e64 v28, s[2:3], v30, v28
	ds_read_b32 v30, v28
	s_waitcnt lgkmcnt(1)
	v_add_u32_e64 v29, s[2:3], v31, v29
	v_lshrrev_b32_e32 v31, 31, v29
	v_add_u32_e64 v29, s[2:3], v29, v31
	v_ashrrev_i32_e32 v29, 1, v29
	s_waitcnt lgkmcnt(0)
	v_sub_u32_e64 v29, s[2:3], v30, v29
	ds_write_b32 v28, v29
	ds_read_b32 v28, v27 offset:44
	s_waitcnt lgkmcnt(0)
	v_mul_lo_u32 v29, v28, s36
	v_mul_lo_u32 v30, v28, s35
	v_lshlrev_b32_e32 v28, 2, v28
	v_lshlrev_b32_e32 v29, 2, v29
	v_lshlrev_b32_e32 v30, 2, v30
	v_add_u32_e64 v29, s[2:3], v19, v29
	v_add_u32_e64 v30, s[2:3], v19, v30
	ds_read_b32 v29, v29
	ds_read_b32 v31, v30
	v_add_u32_e64 v28, s[2:3], v30, v28
	ds_read_b32 v30, v28
	s_waitcnt lgkmcnt(1)
	v_add_u32_e64 v29, s[2:3], v31, v29
	v_lshrrev_b32_e32 v31, 31, v29
	v_add_u32_e64 v29, s[2:3], v29, v31
	v_ashrrev_i32_e32 v29, 1, v29
	s_waitcnt lgkmcnt(0)
	v_sub_u32_e64 v29, s[2:3], v30, v29
	ds_write_b32 v28, v29
	ds_read_b32 v27, v27 offset:44
	s_waitcnt lgkmcnt(0)
	v_mul_lo_u32 v28, v27, s31
	v_mul_lo_u32 v29, v27, s36
	v_lshlrev_b32_e32 v27, 2, v27
	s_add_i32 s31, s31, 8
	v_lshlrev_b32_e32 v28, 2, v28
	v_lshlrev_b32_e32 v29, 2, v29
	v_add_u32_e64 v28, s[2:3], v19, v28
	v_add_u32_e64 v29, s[2:3], v19, v29
	ds_read_b32 v28, v28
	ds_read_b32 v30, v29
	v_add_u32_e64 v27, s[2:3], v29, v27
	ds_read_b32 v29, v27
	s_waitcnt lgkmcnt(1)
	v_add_u32_e64 v28, s[2:3], v30, v28
	v_lshrrev_b32_e32 v30, 31, v28
	v_add_u32_e64 v28, s[2:3], v28, v30
	v_ashrrev_i32_e32 v28, 1, v28
	s_waitcnt lgkmcnt(0)
	v_sub_u32_e64 v28, s[2:3], v29, v28
	ds_write_b32 v27, v28
	s_cbranch_vccnz BB2_205
BB2_206:                                ; %Flow451
                                        ;   in Loop: Header=BB2_177 Depth=1
	v_and_b32_e32 v26, 3, v26
	v_cmp_eq_u32_e32 vcc, 0, v26
	s_and_b64 vcc, exec, vcc
	s_cbranch_vccnz BB2_209
; %bb.207:                              ; %.preheader12.i37.preheader
                                        ;   in Loop: Header=BB2_177 Depth=1
	v_lshlrev_b32_e32 v25, 1, v25
BB2_208:                                ; %.preheader12.i37
                                        ;   Parent Loop BB2_177 Depth=1
                                        ; =>  This Inner Loop Header: Depth=2
	v_mov_b32_e32 v27, 0
	ds_read_b32 v27, v27 offset:44
	v_add_u32_e32 v28, vcc, 2, v25
	v_add_u32_e32 v26, vcc, -1, v26
	v_cmp_ne_u32_e32 vcc, 0, v26
	s_and_b64 vcc, exec, vcc
	s_waitcnt lgkmcnt(0)
	v_mul_lo_u32 v29, v28, v27
	v_mul_lo_u32 v30, v27, v25
	v_mov_b32_e32 v25, v28
	v_lshlrev_b32_e32 v27, 2, v27
	v_lshlrev_b32_e32 v29, 2, v29
	v_lshlrev_b32_e32 v28, 2, v30
	v_add_u32_e64 v29, s[2:3], v19, v29
	v_add_u32_e64 v28, s[2:3], v19, v28
	ds_read_b32 v29, v29
	ds_read_b32 v30, v28
	v_add_u32_e64 v27, s[2:3], v28, v27
	ds_read_b32 v28, v27
	s_waitcnt lgkmcnt(1)
	v_add_u32_e64 v29, s[2:3], v30, v29
	v_lshrrev_b32_e32 v30, 31, v29
	v_add_u32_e64 v29, s[2:3], v29, v30
	v_ashrrev_i32_e32 v29, 1, v29
	s_waitcnt lgkmcnt(0)
	v_sub_u32_e64 v28, s[2:3], v28, v29
	ds_write_b32 v27, v28
	s_cbranch_vccnz BB2_208
BB2_209:                                ; %.loopexit13.i38
                                        ;   in Loop: Header=BB2_177 Depth=1
	v_mov_b32_e32 v25, 0
	ds_read_b32 v26, v25 offset:40
	s_waitcnt lgkmcnt(0)
	v_cmp_gt_i32_e32 vcc, 4, v26
	s_and_b64 vcc, exec, vcc
	s_cbranch_vccnz BB2_216
; %bb.210:                              ;   in Loop: Header=BB2_177 Depth=1
	v_lshrrev_b32_e32 v26, 1, v26
	v_max_u32_e32 v27, 2, v26
	v_add_u32_e32 v26, vcc, -1, v27
	v_add_u32_e32 v27, vcc, -2, v27
	v_cmp_gt_u32_e32 vcc, 3, v27
	s_and_b64 vcc, exec, vcc
	s_cbranch_vccnz BB2_213
; %bb.211:                              ;   in Loop: Header=BB2_177 Depth=1
	v_and_b32_e32 v25, -4, v26
	s_mov_b32 s30, 0
	s_mov_b32 s31, 9
BB2_212:                                ;   Parent Loop BB2_177 Depth=1
                                        ; =>  This Inner Loop Header: Depth=2
	v_mov_b32_e32 v27, 0
	ds_read_b32 v28, v27 offset:44
	s_add_i32 s2, s31, -6
	s_add_i32 s3, s31, -4
	s_add_i32 s33, s31, -2
	s_add_i32 s30, s30, 4
	s_waitcnt lgkmcnt(0)
	v_mul_lo_u32 v29, v28, s2
	v_lshlrev_b32_e32 v31, 1, v28
	v_lshlrev_b32_e32 v28, 2, v28
	v_lshlrev_b32_e32 v30, 2, v29
	v_subrev_u32_e32 v29, vcc, v31, v29
	v_lshlrev_b32_e32 v29, 2, v29
	v_add_u32_e32 v30, vcc, v19, v30
	v_add_u32_e32 v29, vcc, v19, v29
	ds_read_b32 v30, v30
	ds_read_b32 v31, v29
	v_add_u32_e32 v28, vcc, v29, v28
	ds_read_b32 v29, v28
	s_waitcnt lgkmcnt(1)
	v_add_u32_e32 v30, vcc, v30, v31
	v_add_u32_e32 v30, vcc, 2, v30
	v_ashrrev_i32_e32 v31, 31, v30
	v_lshrrev_b32_e32 v31, 30, v31
	v_add_u32_e32 v30, vcc, v30, v31
	v_ashrrev_i32_e32 v30, 2, v30
	s_waitcnt lgkmcnt(0)
	v_add_u32_e32 v29, vcc, v30, v29
	ds_write_b32 v28, v29
	ds_read_b32 v28, v27 offset:44
	s_waitcnt lgkmcnt(0)
	v_mul_lo_u32 v31, v28, s2
	v_mul_lo_u32 v29, v28, s3
	v_lshlrev_b32_e32 v31, 2, v31
	v_lshlrev_b32_e32 v30, 2, v29
	v_add_u32_e32 v30, vcc, v19, v30
	v_add_u32_e32 v31, vcc, v19, v31
	ds_read_b32 v30, v30
	ds_read_b32 v31, v31
	v_subrev_u32_e32 v28, vcc, v28, v29
	v_lshlrev_b32_e32 v28, 2, v28
	v_add_u32_e32 v28, vcc, v19, v28
	ds_read_b32 v29, v28
	s_waitcnt lgkmcnt(1)
	v_add_u32_e32 v30, vcc, v30, v31
	v_add_u32_e32 v30, vcc, 2, v30
	v_ashrrev_i32_e32 v31, 31, v30
	v_lshrrev_b32_e32 v31, 30, v31
	v_add_u32_e32 v30, vcc, v30, v31
	v_ashrrev_i32_e32 v30, 2, v30
	s_waitcnt lgkmcnt(0)
	v_add_u32_e32 v29, vcc, v30, v29
	ds_write_b32 v28, v29
	ds_read_b32 v28, v27 offset:44
	s_waitcnt lgkmcnt(0)
	v_mul_lo_u32 v29, v28, s33
	v_lshlrev_b32_e32 v31, 1, v28
	v_lshlrev_b32_e32 v28, 2, v28
	v_lshlrev_b32_e32 v30, 2, v29
	v_subrev_u32_e32 v29, vcc, v31, v29
	v_lshlrev_b32_e32 v29, 2, v29
	v_add_u32_e32 v30, vcc, v19, v30
	v_add_u32_e32 v29, vcc, v19, v29
	ds_read_b32 v30, v30
	ds_read_b32 v31, v29
	v_add_u32_e32 v28, vcc, v29, v28
	ds_read_b32 v29, v28
	s_waitcnt lgkmcnt(1)
	v_add_u32_e32 v30, vcc, v30, v31
	v_add_u32_e32 v30, vcc, 2, v30
	v_ashrrev_i32_e32 v31, 31, v30
	v_lshrrev_b32_e32 v31, 30, v31
	v_add_u32_e32 v30, vcc, v30, v31
	v_ashrrev_i32_e32 v30, 2, v30
	s_waitcnt lgkmcnt(0)
	v_add_u32_e32 v29, vcc, v30, v29
	ds_write_b32 v28, v29
	ds_read_b32 v27, v27 offset:44
	s_waitcnt lgkmcnt(0)
	v_mul_lo_u32 v28, v27, s31
	v_mul_lo_u32 v30, v27, s33
	s_add_i32 s31, s31, 8
	v_lshlrev_b32_e32 v29, 2, v28
	v_lshlrev_b32_e32 v30, 2, v30
	v_add_u32_e32 v29, vcc, v19, v29
	v_add_u32_e32 v30, vcc, v19, v30
	ds_read_b32 v29, v29
	ds_read_b32 v30, v30
	v_subrev_u32_e32 v27, vcc, v27, v28
	v_lshlrev_b32_e32 v27, 2, v27
	v_add_u32_e64 v27, s[2:3], v19, v27
	ds_read_b32 v28, v27
	s_waitcnt lgkmcnt(1)
	v_add_u32_e64 v29, s[2:3], v29, v30
	v_add_u32_e64 v29, s[2:3], 2, v29
	v_ashrrev_i32_e32 v30, 31, v29
	v_lshrrev_b32_e32 v30, 30, v30
	v_add_u32_e64 v29, s[2:3], v29, v30
	v_cmp_eq_u32_e32 vcc, s30, v25
	v_ashrrev_i32_e32 v29, 2, v29
	s_and_b64 vcc, exec, vcc
	s_waitcnt lgkmcnt(0)
	v_add_u32_e64 v28, s[2:3], v29, v28
	ds_write_b32 v27, v28
	s_cbranch_vccz BB2_212
BB2_213:                                ; %Flow446
                                        ;   in Loop: Header=BB2_177 Depth=1
	v_and_b32_e32 v26, 3, v26
	v_cmp_eq_u32_e32 vcc, 0, v26
	s_and_b64 vcc, exec, vcc
	s_cbranch_vccnz BB2_216
; %bb.214:                              ; %.preheader9.i41.preheader
                                        ;   in Loop: Header=BB2_177 Depth=1
	v_lshlrev_b32_e32 v25, 1, v25
	v_add_u32_e32 v25, vcc, 3, v25
BB2_215:                                ; %.preheader9.i41
                                        ;   Parent Loop BB2_177 Depth=1
                                        ; =>  This Inner Loop Header: Depth=2
	v_mov_b32_e32 v27, 0
	ds_read_b32 v27, v27 offset:44
	v_add_u32_e32 v26, vcc, -1, v26
	s_waitcnt lgkmcnt(0)
	v_mul_lo_u32 v28, v27, v25
	v_lshlrev_b32_e32 v30, 1, v27
	v_lshlrev_b32_e32 v27, 2, v27
	v_add_u32_e32 v25, vcc, 2, v25
	v_lshlrev_b32_e32 v29, 2, v28
	v_subrev_u32_e32 v28, vcc, v30, v28
	v_lshlrev_b32_e32 v28, 2, v28
	v_add_u32_e32 v29, vcc, v19, v29
	v_add_u32_e64 v28, s[2:3], v19, v28
	ds_read_b32 v29, v29
	ds_read_b32 v30, v28
	v_add_u32_e64 v27, s[2:3], v28, v27
	ds_read_b32 v28, v27
	v_cmp_ne_u32_e32 vcc, 0, v26
	s_and_b64 vcc, exec, vcc
	s_waitcnt lgkmcnt(1)
	v_add_u32_e64 v29, s[2:3], v29, v30
	v_add_u32_e64 v29, s[2:3], 2, v29
	v_ashrrev_i32_e32 v30, 31, v29
	v_lshrrev_b32_e32 v30, 30, v30
	v_add_u32_e64 v29, s[2:3], v29, v30
	v_ashrrev_i32_e32 v29, 2, v29
	s_waitcnt lgkmcnt(0)
	v_add_u32_e64 v28, s[2:3], v29, v28
	ds_write_b32 v27, v28
	s_cbranch_vccnz BB2_215
BB2_216:                                ; %Flow456
                                        ;   in Loop: Header=BB2_177 Depth=1
	s_or_b64 exec, exec, s[4:5]
	v_mov_b32_e32 v27, 0
	s_waitcnt lgkmcnt(0)
	s_barrier
	s_waitcnt lgkmcnt(0)
	ds_read2_b32 v[25:26], v27 offset0:1 offset1:9
	s_waitcnt lgkmcnt(0)
	v_ashrrev_i32_e32 v30, 31, v26
	v_add_u32_e32 v28, vcc, v30, v26
	v_xor_b32_e32 v31, v28, v30
	v_cvt_f32_u32_e32 v28, v31
	v_rcp_iflag_f32_e32 v32, v28
	ds_read2_b32 v[28:29], v27 offset0:11 offset1:15
	s_waitcnt lgkmcnt(0)
	v_mul_lo_u32 v27, v28, v25
	v_mul_f32_e32 v25, s28, v32
	v_cvt_u32_f32_e32 v32, v25
	v_lshlrev_b32_e32 v25, 1, v28
	v_add_u32_e32 v28, vcc, -1, v27
	v_ashrrev_i32_e32 v34, 31, v28
	v_mul_lo_u32 v27, v32, v31
	v_mul_hi_u32 v33, v32, v31
	v_xor_b32_e32 v30, v34, v30
	v_sub_u32_e32 v35, vcc, 0, v27
	v_cmp_eq_u32_e64 s[2:3], 0, v33
	v_cndmask_b32_e64 v27, v27, v35, s[2:3]
	v_mul_hi_u32 v27, v27, v32
	v_add_u32_e32 v33, vcc, v34, v28
	v_xor_b32_e32 v33, v33, v34
	v_add_u32_e32 v35, vcc, v27, v32
	v_subrev_u32_e32 v27, vcc, v27, v32
	v_cndmask_b32_e64 v27, v27, v35, s[2:3]
	v_mul_hi_u32 v32, v27, v33
	v_add_u32_e32 v27, vcc, v29, v25
	v_mul_lo_u32 v29, v32, v31
	v_add_u32_e32 v34, vcc, -1, v32
	v_subrev_u32_e32 v35, vcc, v29, v33
	v_cmp_ge_u32_e64 s[4:5], v33, v29
	v_cmp_ge_u32_e64 s[2:3], v35, v31
	v_add_u32_e32 v29, vcc, 1, v32
	s_and_b64 vcc, s[2:3], s[4:5]
	v_cndmask_b32_e32 v29, v32, v29, vcc
	v_cndmask_b32_e64 v29, v34, v29, s[4:5]
	v_xor_b32_e32 v31, v29, v30
	v_sub_u32_e32 v29, vcc, v31, v30
	v_cmp_gt_i32_e32 vcc, 1, v29
	s_and_b64 vcc, exec, vcc
	s_cbranch_vccnz BB2_221
; %bb.217:                              ; %.preheader7.i44.preheader
                                        ;   in Loop: Header=BB2_177 Depth=1
	v_not_b32_e32 v30, v30
	v_add_u32_e32 v30, vcc, v31, v30
	s_mov_b32 s2, 0
	v_mov_b32_e32 v31, v26
	s_branch BB2_219
BB2_218:                                ;   in Loop: Header=BB2_219 Depth=2
	v_mov_b32_e32 v31, 0
	ds_read_b32 v31, v31 offset:36
	s_mov_b64 s[4:5], 0
	s_andn2_b64 vcc, exec, s[4:5]
	s_add_i32 s2, s2, 1
	s_cbranch_vccz BB2_221
BB2_219:                                ; %.preheader7.i44
                                        ;   Parent Loop BB2_177 Depth=1
                                        ; =>  This Inner Loop Header: Depth=2
	s_waitcnt lgkmcnt(0)
	v_mul_lo_u32 v31, s2, v31
	v_add_u32_e32 v31, vcc, v31, v2
	v_add_u32_e32 v32, vcc, v31, v25
	v_add_u32_e32 v31, vcc, v31, v27
	v_lshlrev_b32_e32 v32, 2, v32
	v_lshlrev_b32_e32 v33, 2, v31
	ds_read2_b32 v[31:32], v32 offset0:16 offset1:17
	ds_read_b32 v34, v33 offset:64
	s_waitcnt lgkmcnt(1)
	v_add_u32_e32 v31, vcc, v32, v31
	v_lshrrev_b32_e32 v32, 31, v31
	v_add_u32_e32 v31, vcc, v31, v32
	v_ashrrev_i32_e32 v31, 1, v31
	s_waitcnt lgkmcnt(0)
	v_sub_u32_e32 v31, vcc, v34, v31
	v_cmp_eq_u32_e32 vcc, s2, v30
	s_and_b64 vcc, exec, vcc
	ds_write_b32 v33, v31 offset:64
	s_cbranch_vccz BB2_218
; %bb.220:                              ;   in Loop: Header=BB2_219 Depth=2
	s_mov_b64 s[4:5], -1
                                        ; implicit-def: $vgpr31
	s_andn2_b64 vcc, exec, s[4:5]
	s_add_i32 s2, s2, 1
	s_cbranch_vccnz BB2_219
BB2_221:                                ; %Flow442
                                        ;   in Loop: Header=BB2_177 Depth=1
	v_mul_lo_u32 v26, v29, v26
	v_subrev_u32_e32 v29, vcc, v26, v28
	v_subrev_u32_e32 v26, vcc, v29, v28
	v_ashrrev_i32_e32 v30, 31, v29
	v_cmp_lt_u64_e32 vcc, v[2:3], v[29:30]
	s_and_saveexec_b64 s[2:3], vcc
	s_cbranch_execz BB2_223
; %bb.222:                              ;   in Loop: Header=BB2_177 Depth=1
	v_add_u32_e32 v26, vcc, v26, v2
	v_add_u32_e32 v25, vcc, v26, v25
	v_add_u32_e32 v26, vcc, v26, v27
	v_lshlrev_b32_e32 v25, 2, v25
	v_lshlrev_b32_e32 v27, 2, v26
	ds_read2_b32 v[25:26], v25 offset0:16 offset1:17
	ds_read_b32 v28, v27 offset:64
	s_waitcnt lgkmcnt(1)
	v_add_u32_e32 v25, vcc, v26, v25
	v_lshrrev_b32_e32 v26, 31, v25
	v_add_u32_e32 v25, vcc, v25, v26
	v_ashrrev_i32_e32 v25, 1, v25
	s_waitcnt lgkmcnt(0)
	v_sub_u32_e32 v25, vcc, v28, v25
	ds_write_b32 v27, v25 offset:64
BB2_223:                                ;   in Loop: Header=BB2_177 Depth=1
	s_or_b64 exec, exec, s[2:3]
	v_mov_b32_e32 v27, 0
	s_waitcnt lgkmcnt(0)
	s_barrier
	s_waitcnt lgkmcnt(0)
	ds_read2_b32 v[25:26], v27 offset0:1 offset1:9
	s_waitcnt lgkmcnt(0)
	v_ashrrev_i32_e32 v30, 31, v26
	v_add_u32_e32 v28, vcc, v30, v26
	v_xor_b32_e32 v31, v28, v30
	v_cvt_f32_u32_e32 v28, v31
	v_rcp_iflag_f32_e32 v32, v28
	ds_read2_b32 v[28:29], v27 offset0:11 offset1:15
	v_mul_f32_e32 v27, s28, v32
	v_cvt_u32_f32_e32 v32, v27
	s_waitcnt lgkmcnt(0)
	v_mul_lo_u32 v25, v28, v25
	v_lshlrev_b32_e32 v28, 1, v28
	v_mul_hi_u32 v33, v32, v31
	v_add_u32_e32 v27, vcc, -1, v25
	v_mul_lo_u32 v25, v32, v31
	v_ashrrev_i32_e32 v34, 31, v27
	v_cmp_eq_u32_e64 s[2:3], 0, v33
	v_add_u32_e32 v33, vcc, v34, v27
	v_sub_u32_e32 v35, vcc, 0, v25
	v_cndmask_b32_e64 v25, v25, v35, s[2:3]
	v_mul_hi_u32 v25, v25, v32
	v_xor_b32_e32 v33, v33, v34
	v_xor_b32_e32 v30, v34, v30
	v_add_u32_e32 v35, vcc, v25, v32
	v_subrev_u32_e32 v25, vcc, v25, v32
	v_cndmask_b32_e64 v25, v25, v35, s[2:3]
	v_mul_hi_u32 v32, v25, v33
	v_add_u32_e32 v25, vcc, v29, v28
	v_mul_lo_u32 v29, v32, v31
	v_add_u32_e32 v34, vcc, -1, v32
	v_subrev_u32_e32 v35, vcc, v29, v33
	v_cmp_ge_u32_e64 s[4:5], v33, v29
	v_cmp_ge_u32_e64 s[2:3], v35, v31
	v_add_u32_e32 v29, vcc, 1, v32
	s_and_b64 vcc, s[2:3], s[4:5]
	v_cndmask_b32_e32 v29, v32, v29, vcc
	v_cndmask_b32_e64 v29, v34, v29, s[4:5]
	v_xor_b32_e32 v31, v29, v30
	v_sub_u32_e32 v29, vcc, v31, v30
	v_cmp_gt_i32_e32 vcc, 1, v29
	s_and_b64 vcc, exec, vcc
	s_cbranch_vccnz BB2_228
; %bb.224:                              ; %.preheader5.i46.preheader
                                        ;   in Loop: Header=BB2_177 Depth=1
	v_not_b32_e32 v30, v30
	v_add_u32_e32 v30, vcc, v31, v30
	s_mov_b32 s2, 0
	v_mov_b32_e32 v31, v26
	s_branch BB2_226
BB2_225:                                ;   in Loop: Header=BB2_226 Depth=2
	v_mov_b32_e32 v31, 0
	ds_read_b32 v31, v31 offset:36
	s_mov_b64 s[4:5], 0
	s_andn2_b64 vcc, exec, s[4:5]
	s_add_i32 s2, s2, 1
	s_cbranch_vccz BB2_228
BB2_226:                                ; %.preheader5.i46
                                        ;   Parent Loop BB2_177 Depth=1
                                        ; =>  This Inner Loop Header: Depth=2
	s_waitcnt lgkmcnt(0)
	v_mul_lo_u32 v31, s2, v31
	v_add_u32_e32 v31, vcc, v31, v2
	v_add_u32_e32 v32, vcc, v31, v25
	v_add_u32_e32 v31, vcc, v31, v28
	v_lshlrev_b32_e32 v32, 2, v32
	v_lshlrev_b32_e32 v33, 2, v31
	ds_read2_b32 v[31:32], v32 offset0:16 offset1:17
	ds_read_b32 v34, v33 offset:68
	s_waitcnt lgkmcnt(1)
	v_add_u32_e32 v31, vcc, v31, v32
	v_add_u32_e32 v31, vcc, 2, v31
	v_ashrrev_i32_e32 v32, 31, v31
	v_lshrrev_b32_e32 v32, 30, v32
	v_add_u32_e32 v31, vcc, v31, v32
	v_ashrrev_i32_e32 v31, 2, v31
	s_waitcnt lgkmcnt(0)
	v_add_u32_e32 v31, vcc, v31, v34
	v_cmp_eq_u32_e32 vcc, s2, v30
	s_and_b64 vcc, exec, vcc
	ds_write_b32 v33, v31 offset:68
	s_cbranch_vccz BB2_225
; %bb.227:                              ;   in Loop: Header=BB2_226 Depth=2
	s_mov_b64 s[4:5], -1
                                        ; implicit-def: $vgpr31
	s_andn2_b64 vcc, exec, s[4:5]
	s_add_i32 s2, s2, 1
	s_cbranch_vccnz BB2_226
BB2_228:                                ; %Flow439
                                        ;   in Loop: Header=BB2_177 Depth=1
	v_mul_lo_u32 v29, v29, v26
	v_or_b32_e32 v26, 1, v28
	v_subrev_u32_e32 v28, vcc, v29, v27
	v_subrev_u32_e32 v27, vcc, v28, v27
	v_ashrrev_i32_e32 v29, 31, v28
	v_cmp_lt_u64_e32 vcc, v[2:3], v[28:29]
	s_and_saveexec_b64 s[2:3], vcc
	s_cbranch_execz BB2_230
; %bb.229:                              ;   in Loop: Header=BB2_177 Depth=1
	v_add_u32_e32 v27, vcc, v27, v2
	v_add_u32_e32 v25, vcc, v27, v25
	v_add_u32_e32 v26, vcc, v27, v26
	v_lshlrev_b32_e32 v25, 2, v25
	v_lshlrev_b32_e32 v27, 2, v26
	ds_read2_b32 v[25:26], v25 offset0:16 offset1:17
	ds_read_b32 v28, v27 offset:64
	s_waitcnt lgkmcnt(1)
	v_add_u32_e32 v25, vcc, v25, v26
	v_add_u32_e32 v25, vcc, 2, v25
	v_ashrrev_i32_e32 v26, 31, v25
	v_lshrrev_b32_e32 v26, 30, v26
	v_add_u32_e32 v25, vcc, v25, v26
	v_ashrrev_i32_e32 v25, 2, v25
	s_waitcnt lgkmcnt(0)
	v_add_u32_e32 v25, vcc, v25, v28
	ds_write_b32 v27, v25 offset:64
BB2_230:                                ;   in Loop: Header=BB2_177 Depth=1
	s_or_b64 exec, exec, s[2:3]
	v_mov_b32_e32 v25, 0
	s_waitcnt lgkmcnt(0)
	s_barrier
	s_waitcnt lgkmcnt(0)
	ds_read_b32 v26, v25 offset:4
	s_waitcnt lgkmcnt(0)
	v_cmp_gt_i32_e32 vcc, 1, v26
	s_and_b64 vcc, exec, vcc
	s_cbranch_vccnz BB2_176
; %bb.231:                              ;   in Loop: Header=BB2_177 Depth=1
	v_add_u32_e32 v26, vcc, -1, v26
	v_lshrrev_b32_e32 v27, 1, v26
	v_add_u32_e32 v27, vcc, 1, v27
	v_cmp_gt_u32_e32 vcc, 6, v26
	ds_read_b32 v26, v25 offset:44
	s_and_b64 vcc, exec, vcc
	s_cbranch_vccnz BB2_250
; %bb.232:                              ;   in Loop: Header=BB2_177 Depth=1
	s_waitcnt lgkmcnt(0)
	v_mul_lo_u32 v25, v26, 9
	v_mul_lo_u32 v31, v26, 7
	v_mul_lo_u32 v32, v26, 6
	v_mul_lo_u32 v35, v26, 12
	v_add_u32_e32 v25, vcc, v22, v25
	v_lshlrev_b32_e32 v29, 2, v25
	v_lshlrev_b32_e32 v25, 3, v26
	v_add_u32_e32 v25, vcc, v22, v25
	v_lshlrev_b32_e32 v30, 2, v25
	v_add_u32_e32 v25, vcc, v22, v31
	v_lshlrev_b32_e32 v31, 2, v25
	v_mul_lo_u32 v25, v26, 5
	v_add_u32_e32 v32, vcc, v22, v32
	s_mov_b32 s30, 2
	v_lshlrev_b32_e32 v28, 5, v26
	v_add_u32_e32 v25, vcc, v22, v25
	v_lshlrev_b32_e32 v33, 2, v25
	v_lshlrev_b32_e32 v25, 2, v26
	v_add_u32_e32 v25, vcc, v22, v25
	v_lshlrev_b32_e32 v34, 2, v25
	v_lshlrev_b32_e32 v25, 1, v26
	v_add_u32_e32 v25, vcc, v22, v25
	v_lshlrev_b32_e32 v36, 2, v25
	v_and_b32_e32 v25, s27, v27
	v_lshlrev_b32_e32 v32, 2, v32
	v_add_u32_e32 v35, vcc, v23, v35
	v_sub_u32_e32 v37, vcc, 0, v25
	v_mov_b32_e32 v38, 0
	s_branch BB2_235
BB2_233:                                ; %Flow430
                                        ;   in Loop: Header=BB2_235 Depth=2
	s_or_b64 exec, exec, s[4:5]
BB2_234:                                ;   in Loop: Header=BB2_235 Depth=2
	s_or_b64 exec, exec, s[34:35]
	v_cndmask_b32_e64 v14, 0, 1, vcc
	v_add_u32_e32 v38, vcc, v38, v28
	v_add_u32_e32 v37, vcc, 4, v37
	v_cmp_eq_u32_e32 vcc, 0, v37
	s_and_b64 vcc, exec, vcc
	s_add_i32 s30, s30, 8
	s_cbranch_vccnz BB2_251
BB2_235:                                ;   Parent Loop BB2_177 Depth=1
                                        ; =>  This Inner Loop Header: Depth=2
	v_and_b32_e32 v14, s29, v14
	v_cmp_eq_u16_e64 s[2:3], 0, v14
	v_cmp_ne_u32_e64 s[4:5], v24, v15
	v_cmp_ne_u16_e32 vcc, 0, v14
	s_or_b64 s[2:3], s[2:3], s[4:5]
	v_mov_b32_e32 v39, v15
	s_and_saveexec_b64 s[4:5], s[2:3]
	s_cbranch_execz BB2_239
; %bb.236:                              ;   in Loop: Header=BB2_235 Depth=2
	v_add_u32_e64 v14, s[2:3], v38, v36
	v_ashrrev_i32_e32 v25, 31, v24
	ds_read_b32 v14, v14 offset:64
	v_lshlrev_b64 v[40:41], 2, v[24:25]
	v_mov_b32_e32 v25, s11
	v_add_u32_e64 v40, s[2:3], s10, v40
	v_addc_u32_e64 v41, s[2:3], v25, v41, s[2:3]
	s_waitcnt lgkmcnt(0)
	flat_store_dword v[40:41], v14
	v_add_u32_e64 v14, s[2:3], v24, v13
	v_cmp_ne_u32_e64 s[2:3], v14, v15
	s_xor_b64 s[34:35], vcc, -1
	s_or_b64 s[2:3], s[34:35], s[2:3]
	v_mov_b32_e32 v39, v15
	s_and_saveexec_b64 s[34:35], s[2:3]
	s_cbranch_execz BB2_238
; %bb.237:                              ;   in Loop: Header=BB2_235 Depth=2
	v_add_u32_e64 v24, s[2:3], v38, v35
	ds_read_b32 v39, v24 offset:64
	v_ashrrev_i32_e32 v25, 31, v13
	v_mov_b32_e32 v24, v13
	v_lshlrev_b64 v[24:25], 2, v[24:25]
	v_add_u32_e64 v24, s[2:3], v40, v24
	v_addc_u32_e64 v25, s[2:3], v41, v25, s[2:3]
	s_waitcnt lgkmcnt(0)
	flat_store_dword v[24:25], v39
	v_add_u32_e64 v39, s[2:3], v14, v16
BB2_238:                                ; %Flow433
                                        ;   in Loop: Header=BB2_235 Depth=2
	s_or_b64 exec, exec, s[34:35]
BB2_239:                                ;   in Loop: Header=BB2_235 Depth=2
	s_or_b64 exec, exec, s[4:5]
	v_cmp_ne_u32_e64 s[2:3], v39, v15
	s_xor_b64 s[4:5], vcc, -1
	s_or_b64 s[2:3], s[4:5], s[2:3]
	v_mov_b32_e32 v24, v15
	s_and_saveexec_b64 s[34:35], s[2:3]
	s_cbranch_execz BB2_243
; %bb.240:                              ;   in Loop: Header=BB2_235 Depth=2
	v_add_u32_e64 v14, s[2:3], v38, v34
	v_ashrrev_i32_e32 v40, 31, v39
	ds_read_b32 v14, v14 offset:64
	v_lshlrev_b64 v[24:25], 2, v[39:40]
	v_mov_b32_e32 v41, s11
	v_add_u32_e64 v40, s[2:3], s10, v24
	v_addc_u32_e64 v41, s[2:3], v41, v25, s[2:3]
	s_waitcnt lgkmcnt(0)
	flat_store_dword v[40:41], v14
	v_add_u32_e64 v14, s[2:3], v39, v13
	v_cmp_ne_u32_e64 s[2:3], v14, v15
	s_or_b64 s[2:3], s[4:5], s[2:3]
	v_mov_b32_e32 v24, v15
	s_and_saveexec_b64 s[36:37], s[2:3]
	s_cbranch_execz BB2_242
; %bb.241:                              ;   in Loop: Header=BB2_235 Depth=2
	v_add_u32_e64 v24, s[2:3], v38, v33
	ds_read_b32 v39, v24 offset:64
	v_ashrrev_i32_e32 v25, 31, v13
	v_mov_b32_e32 v24, v13
	v_lshlrev_b64 v[24:25], 2, v[24:25]
	v_add_u32_e64 v24, s[2:3], v40, v24
	v_addc_u32_e64 v25, s[2:3], v41, v25, s[2:3]
	s_waitcnt lgkmcnt(0)
	flat_store_dword v[24:25], v39
	v_add_u32_e64 v24, s[2:3], v14, v16
BB2_242:                                ; %Flow432
                                        ;   in Loop: Header=BB2_235 Depth=2
	s_or_b64 exec, exec, s[36:37]
BB2_243:                                ;   in Loop: Header=BB2_235 Depth=2
	s_or_b64 exec, exec, s[34:35]
	v_cmp_ne_u32_e64 s[2:3], v24, v15
	s_or_b64 s[2:3], s[4:5], s[2:3]
	v_mov_b32_e32 v39, v15
	s_and_saveexec_b64 s[34:35], s[2:3]
	s_cbranch_execz BB2_247
; %bb.244:                              ;   in Loop: Header=BB2_235 Depth=2
	v_add_u32_e64 v14, s[2:3], v38, v32
	v_ashrrev_i32_e32 v25, 31, v24
	ds_read_b32 v14, v14 offset:64
	v_lshlrev_b64 v[40:41], 2, v[24:25]
	v_mov_b32_e32 v25, s11
	v_add_u32_e64 v40, s[2:3], s10, v40
	v_addc_u32_e64 v41, s[2:3], v25, v41, s[2:3]
	s_waitcnt lgkmcnt(0)
	flat_store_dword v[40:41], v14
	v_add_u32_e64 v14, s[2:3], v24, v13
	v_cmp_ne_u32_e64 s[2:3], v14, v15
	s_or_b64 s[2:3], s[4:5], s[2:3]
	v_mov_b32_e32 v39, v15
	s_and_saveexec_b64 s[36:37], s[2:3]
	s_cbranch_execz BB2_246
; %bb.245:                              ;   in Loop: Header=BB2_235 Depth=2
	v_add_u32_e64 v24, s[2:3], v38, v31
	ds_read_b32 v39, v24 offset:64
	v_ashrrev_i32_e32 v25, 31, v13
	v_mov_b32_e32 v24, v13
	v_lshlrev_b64 v[24:25], 2, v[24:25]
	v_add_u32_e64 v24, s[2:3], v40, v24
	v_addc_u32_e64 v25, s[2:3], v41, v25, s[2:3]
	s_waitcnt lgkmcnt(0)
	flat_store_dword v[24:25], v39
	v_add_u32_e64 v39, s[2:3], v14, v16
BB2_246:                                ; %Flow431
                                        ;   in Loop: Header=BB2_235 Depth=2
	s_or_b64 exec, exec, s[36:37]
BB2_247:                                ;   in Loop: Header=BB2_235 Depth=2
	s_or_b64 exec, exec, s[34:35]
	v_cmp_ne_u32_e64 s[2:3], v39, v15
	s_or_b64 s[2:3], s[4:5], s[2:3]
	v_mov_b32_e32 v24, v15
	s_and_saveexec_b64 s[34:35], s[2:3]
	s_cbranch_execz BB2_234
; %bb.248:                              ;   in Loop: Header=BB2_235 Depth=2
	v_add_u32_e64 v14, s[2:3], v38, v30
	v_ashrrev_i32_e32 v40, 31, v39
	ds_read_b32 v14, v14 offset:64
	v_lshlrev_b64 v[24:25], 2, v[39:40]
	v_mov_b32_e32 v41, s11
	v_add_u32_e64 v40, s[2:3], s10, v24
	v_addc_u32_e64 v41, s[2:3], v41, v25, s[2:3]
	s_waitcnt lgkmcnt(0)
	flat_store_dword v[40:41], v14
	v_add_u32_e64 v14, s[2:3], v39, v13
	v_cmp_ne_u32_e64 s[2:3], v14, v15
	s_or_b64 s[2:3], s[4:5], s[2:3]
	v_mov_b32_e32 v24, v15
	s_and_saveexec_b64 s[4:5], s[2:3]
	s_cbranch_execz BB2_233
; %bb.249:                              ;   in Loop: Header=BB2_235 Depth=2
	v_add_u32_e64 v24, s[2:3], v38, v29
	ds_read_b32 v39, v24 offset:64
	v_ashrrev_i32_e32 v25, 31, v13
	v_mov_b32_e32 v24, v13
	v_lshlrev_b64 v[24:25], 2, v[24:25]
	v_add_u32_e64 v24, s[2:3], v40, v24
	v_addc_u32_e64 v25, s[2:3], v41, v25, s[2:3]
	s_waitcnt lgkmcnt(0)
	flat_store_dword v[24:25], v39
	v_add_u32_e64 v24, s[2:3], v14, v16
	s_branch BB2_233
BB2_250:                                ;   in Loop: Header=BB2_177 Depth=1
	s_mov_b32 s30, 2
BB2_251:                                ; %Flow435
                                        ;   in Loop: Header=BB2_177 Depth=1
	v_and_b32_e32 v27, 3, v27
	v_cmp_eq_u32_e32 vcc, 0, v27
	s_and_b64 vcc, exec, vcc
	s_cbranch_vccnz BB2_176
; %bb.252:                              ; %.preheader.i48.preheader
                                        ;   in Loop: Header=BB2_177 Depth=1
	s_waitcnt lgkmcnt(0)
	v_mul_lo_u32 v25, s30, v26
	v_lshlrev_b32_e32 v29, 3, v26
	v_add_u32_e32 v25, vcc, v22, v25
	v_lshlrev_b32_e32 v28, 2, v25
	s_branch BB2_255
BB2_253:                                ; %Flow427
                                        ;   in Loop: Header=BB2_255 Depth=2
	s_or_b64 exec, exec, s[34:35]
BB2_254:                                ;   in Loop: Header=BB2_255 Depth=2
	s_or_b64 exec, exec, s[4:5]
	v_cndmask_b32_e64 v14, 0, 1, vcc
	v_add_u32_e32 v27, vcc, -1, v27
	v_add_u32_e32 v28, vcc, v28, v29
	v_cmp_ne_u32_e32 vcc, 0, v27
	s_add_i32 s30, s30, 2
	s_and_b64 vcc, exec, vcc
	v_mov_b32_e32 v24, v25
	s_cbranch_vccz BB2_175
BB2_255:                                ; %.preheader.i48
                                        ;   Parent Loop BB2_177 Depth=1
                                        ; =>  This Inner Loop Header: Depth=2
	v_and_b32_e32 v14, s29, v14
	v_cmp_eq_u16_e64 s[2:3], 0, v14
	v_cmp_ne_u32_e64 s[4:5], v24, v15
	v_cmp_ne_u16_e32 vcc, 0, v14
	s_or_b64 s[2:3], s[2:3], s[4:5]
	v_mov_b32_e32 v25, v15
	s_and_saveexec_b64 s[4:5], s[2:3]
	s_cbranch_execz BB2_254
; %bb.256:                              ;   in Loop: Header=BB2_255 Depth=2
	v_ashrrev_i32_e32 v25, 31, v24
	v_lshlrev_b64 v[30:31], 2, v[24:25]
	ds_read_b32 v14, v28 offset:64
	v_mov_b32_e32 v25, s11
	v_add_u32_e64 v30, s[2:3], s10, v30
	v_addc_u32_e64 v31, s[2:3], v25, v31, s[2:3]
	v_add_u32_e64 v24, s[2:3], v24, v13
	v_cmp_ne_u32_e64 s[2:3], v24, v15
	s_xor_b64 s[34:35], vcc, -1
	s_or_b64 s[2:3], s[34:35], s[2:3]
	v_mov_b32_e32 v25, v15
	s_waitcnt lgkmcnt(0)
	flat_store_dword v[30:31], v14
	s_and_saveexec_b64 s[34:35], s[2:3]
	s_cbranch_execz BB2_253
; %bb.257:                              ;   in Loop: Header=BB2_255 Depth=2
	s_or_b32 s2, s30, 1
	v_mul_lo_u32 v25, s2, v26
	v_ashrrev_i32_e32 v14, 31, v13
	v_lshlrev_b64 v[32:33], 2, v[13:14]
	v_add_u32_e64 v14, s[2:3], v25, v17
	v_lshlrev_b32_e32 v14, 2, v14
	ds_read_b32 v14, v14 offset:64
	v_add_u32_e64 v30, s[2:3], v30, v32
	v_addc_u32_e64 v31, s[2:3], v31, v33, s[2:3]
	v_add_u32_e64 v25, s[2:3], v24, v16
	s_waitcnt lgkmcnt(0)
	flat_store_dword v[30:31], v14
	s_branch BB2_253
BB2_258:                                ; %Flow471
	s_mov_b64 s[0:1], 0
BB2_259:                                ; %Flow541
	s_and_b64 vcc, exec, s[0:1]
	s_cbranch_vccz BB2_385
; %bb.260:
	v_add_u32_e32 v1, vcc, s23, v0
	v_cmp_ne_u32_e64 s[2:3], s7, 0
	s_and_b64 vcc, exec, s[2:3]
	s_cbranch_vccz BB2_266
; %bb.261:
	v_cmp_le_i32_e32 vcc, s12, v1
                                        ; implicit-def: $vgpr2
	s_and_saveexec_b64 s[0:1], vcc
	s_xor_b64 s[0:1], exec, s[0:1]
; %bb.262:
	s_lshl_b32 s4, s12, 1
	v_sub_u32_e32 v2, vcc, s4, v1
	v_add_u32_e32 v2, vcc, -2, v2
; %bb.263:                              ; %Flow536
	s_or_saveexec_b64 s[0:1], s[0:1]
	s_xor_b64 exec, exec, s[0:1]
; %bb.264:
	v_ashrrev_i32_e32 v2, 31, v1
	v_add_u32_e32 v3, vcc, v1, v2
	v_xor_b32_e32 v2, v3, v2
; %bb.265:
	s_or_b64 exec, exec, s[0:1]
	s_add_i32 s0, s21, -2
	s_mul_i32 s0, s0, s12
	v_add_u32_e32 v3, vcc, s0, v2
	s_add_i32 s0, s12, s12
	s_waitcnt vmcnt(0) lgkmcnt(0)
	v_ashrrev_i32_e32 v4, 31, v3
	v_add_u32_e32 v2, vcc, s0, v3
	v_lshlrev_b64 v[3:4], 2, v[3:4]
	v_mov_b32_e32 v5, s9
	v_add_u32_e32 v3, vcc, s8, v3
	s_ashr_i32 s1, s12, 31
	s_mov_b32 s0, s16
	v_addc_u32_e32 v4, vcc, v5, v4, vcc
	s_lshl_b64 s[0:1], s[0:1], 2
	v_mov_b32_e32 v5, s1
	v_add_u32_e32 v6, vcc, s0, v3
	v_addc_u32_e32 v7, vcc, v4, v5, vcc
	v_add_u32_e32 v8, vcc, s0, v6
	v_addc_u32_e32 v9, vcc, v7, v5, vcc
	flat_load_dword v5, v[3:4]
	flat_load_dword v4, v[6:7]
	flat_load_dword v3, v[8:9]
	s_mov_b64 s[0:1], 0
	s_and_b64 vcc, exec, s[0:1]
	s_cbranch_vccnz BB2_267
	s_branch BB2_276
BB2_266:
                                        ; implicit-def: $vgpr3
	s_waitcnt vmcnt(0) lgkmcnt(0)
                                        ; implicit-def: $vgpr4
                                        ; implicit-def: $vgpr2
                                        ; implicit-def: $vgpr5
	s_cbranch_execz BB2_276
BB2_267:
	v_cmp_le_i32_e32 vcc, s12, v1
                                        ; implicit-def: $vgpr2
	s_and_saveexec_b64 s[0:1], vcc
	s_xor_b64 s[4:5], exec, s[0:1]
; %bb.268:
	s_lshl_b32 s0, s12, 1
	v_sub_u32_e64 v2, s[0:1], s0, v1
	v_add_u32_e64 v2, s[0:1], -2, v2
; %bb.269:                              ; %Flow538
	s_or_saveexec_b64 s[4:5], s[4:5]
	s_xor_b64 exec, exec, s[4:5]
	s_cbranch_execz BB2_271
; %bb.270:
	v_ashrrev_i32_e32 v2, 31, v1
	s_waitcnt vmcnt(0) lgkmcnt(0)
	v_add_u32_e64 v3, s[0:1], v1, v2
	v_xor_b32_e32 v2, v3, v2
BB2_271:
	s_or_b64 exec, exec, s[4:5]
	s_mul_i32 s4, s21, s12
	v_add_u32_e64 v2, s[0:1], s4, v2
	s_waitcnt vmcnt(0) lgkmcnt(0)
	v_ashrrev_i32_e32 v3, 31, v2
	v_lshlrev_b64 v[2:3], 2, v[2:3]
	v_mov_b32_e32 v4, s9
	v_add_u32_e64 v2, s[0:1], s8, v2
	v_addc_u32_e64 v3, s[0:1], v4, v3, s[0:1]
	s_ashr_i32 s1, s12, 31
	s_mov_b32 s0, s16
	s_lshl_b64 s[24:25], s[0:1], 2
	v_mov_b32_e32 v5, s25
	v_add_u32_e64 v4, s[0:1], s24, v2
	v_addc_u32_e64 v5, s[0:1], v3, v5, s[0:1]
	v_mov_b32_e32 v7, s25
	v_add_u32_e64 v6, s[0:1], s24, v4
	v_addc_u32_e64 v7, s[0:1], v5, v7, s[0:1]
	flat_load_dword v3, v[2:3]
	flat_load_dword v4, v[4:5]
	flat_load_dword v5, v[6:7]
                                        ; implicit-def: $vgpr2
	s_and_saveexec_b64 s[0:1], vcc
	s_xor_b64 s[0:1], exec, s[0:1]
; %bb.272:
	s_lshl_b32 s5, s12, 1
	v_sub_u32_e32 v2, vcc, s5, v1
	v_add_u32_e32 v2, vcc, -2, v2
; %bb.273:                              ; %Flow537
	s_or_saveexec_b64 s[0:1], s[0:1]
	s_xor_b64 exec, exec, s[0:1]
; %bb.274:
	v_ashrrev_i32_e32 v2, 31, v1
	v_add_u32_e32 v1, vcc, v1, v2
	v_xor_b32_e32 v2, v1, v2
; %bb.275:
	s_or_b64 exec, exec, s[0:1]
	v_add_u32_e32 v2, vcc, s4, v2
BB2_276:                                ; %Flow540
	v_add_u32_e32 v6, vcc, 2, v0
	v_and_b32_e32 v1, 1, v6
	v_mul_lo_u32 v7, s20, v1
	v_mov_b32_e32 v1, 0
	v_lshrrev_b32_e32 v6, 1, v6
	v_cmp_gt_u32_e64 s[0:1], 3, v0
	v_add_u32_e32 v8, vcc, v7, v6
	v_mov_b32_e32 v9, v1
	v_mov_b32_e32 v10, v1
	v_mov_b32_e32 v11, v1
	v_mov_b32_e32 v13, v1
	v_mov_b32_e32 v6, v1
	v_mov_b32_e32 v7, v1
	s_and_saveexec_b64 s[4:5], s[0:1]
	s_cbranch_execz BB2_294
; %bb.277:
	v_mov_b32_e32 v6, s15
	v_cmp_eq_u32_e32 vcc, 0, v0
	v_cndmask_b32_e32 v6, -3, v6, vcc
	v_add_u32_e32 v7, vcc, v6, v0
	v_add_u32_e32 v12, vcc, s23, v7
	s_and_b64 vcc, exec, s[2:3]
	s_cbranch_vccz BB2_283
; %bb.278:
	v_cmp_le_i32_e32 vcc, s12, v12
                                        ; implicit-def: $vgpr6
	s_and_saveexec_b64 s[2:3], vcc
	s_xor_b64 s[2:3], exec, s[2:3]
; %bb.279:
	s_lshl_b32 s7, s12, 1
	v_sub_u32_e32 v6, vcc, s7, v12
	v_add_u32_e32 v6, vcc, -2, v6
; %bb.280:                              ; %Flow531
	s_or_saveexec_b64 s[2:3], s[2:3]
	s_xor_b64 exec, exec, s[2:3]
; %bb.281:
	v_ashrrev_i32_e32 v6, 31, v12
	v_add_u32_e32 v9, vcc, v12, v6
	v_xor_b32_e32 v6, v9, v6
; %bb.282:
	s_or_b64 exec, exec, s[2:3]
	s_add_i32 s2, s21, -2
	s_mul_i32 s2, s2, s12
	v_add_u32_e32 v9, vcc, s2, v6
	s_add_i32 s2, s12, s12
	v_ashrrev_i32_e32 v10, 31, v9
	v_add_u32_e32 v6, vcc, s2, v9
	v_lshlrev_b64 v[9:10], 2, v[9:10]
	v_mov_b32_e32 v11, s9
	v_add_u32_e32 v9, vcc, s8, v9
	s_ashr_i32 s3, s12, 31
	s_mov_b32 s2, s16
	v_addc_u32_e32 v10, vcc, v11, v10, vcc
	s_lshl_b64 s[2:3], s[2:3], 2
	v_mov_b32_e32 v11, s3
	v_add_u32_e32 v13, vcc, s2, v9
	v_addc_u32_e32 v14, vcc, v10, v11, vcc
	v_add_u32_e32 v15, vcc, s2, v13
	v_addc_u32_e32 v16, vcc, v14, v11, vcc
	flat_load_dword v11, v[9:10]
	flat_load_dword v10, v[13:14]
	flat_load_dword v9, v[15:16]
	s_mov_b64 s[2:3], 0
	s_and_b64 vcc, exec, s[2:3]
	s_cbranch_vccnz BB2_284
	s_branch BB2_293
BB2_283:
                                        ; implicit-def: $vgpr6
                                        ; implicit-def: $vgpr11
                                        ; implicit-def: $vgpr10
                                        ; implicit-def: $vgpr9
	s_cbranch_execz BB2_293
BB2_284:
	v_cmp_le_i32_e32 vcc, s12, v12
                                        ; implicit-def: $vgpr6
	s_and_saveexec_b64 s[2:3], vcc
	s_xor_b64 s[24:25], exec, s[2:3]
; %bb.285:
	s_lshl_b32 s2, s12, 1
	v_sub_u32_e64 v6, s[2:3], s2, v12
	v_add_u32_e64 v6, s[2:3], -2, v6
; %bb.286:                              ; %Flow533
	s_or_saveexec_b64 s[24:25], s[24:25]
	s_xor_b64 exec, exec, s[24:25]
	s_cbranch_execz BB2_288
; %bb.287:
	v_ashrrev_i32_e32 v6, 31, v12
	s_waitcnt vmcnt(0) lgkmcnt(0)
	v_add_u32_e64 v9, s[2:3], v12, v6
	v_xor_b32_e32 v6, v9, v6
BB2_288:
	s_or_b64 exec, exec, s[24:25]
	s_mul_i32 s7, s21, s12
	s_waitcnt vmcnt(0) lgkmcnt(0)
	v_add_u32_e64 v9, s[2:3], s7, v6
	v_ashrrev_i32_e32 v10, 31, v9
	v_lshlrev_b64 v[9:10], 2, v[9:10]
	v_mov_b32_e32 v6, s9
	v_add_u32_e64 v9, s[2:3], s8, v9
	v_addc_u32_e64 v10, s[2:3], v6, v10, s[2:3]
	s_ashr_i32 s3, s12, 31
	s_mov_b32 s2, s16
	s_lshl_b64 s[24:25], s[2:3], 2
	v_mov_b32_e32 v6, s25
	v_add_u32_e64 v13, s[2:3], s24, v9
	v_addc_u32_e64 v14, s[2:3], v10, v6, s[2:3]
	v_add_u32_e64 v15, s[2:3], s24, v13
	v_addc_u32_e64 v16, s[2:3], v14, v6, s[2:3]
	flat_load_dword v9, v[9:10]
	flat_load_dword v10, v[13:14]
	flat_load_dword v11, v[15:16]
                                        ; implicit-def: $vgpr6
	s_and_saveexec_b64 s[2:3], vcc
	s_xor_b64 s[2:3], exec, s[2:3]
; %bb.289:
	s_lshl_b32 s17, s12, 1
	v_sub_u32_e32 v6, vcc, s17, v12
	v_add_u32_e32 v6, vcc, -2, v6
; %bb.290:                              ; %Flow532
	s_or_saveexec_b64 s[2:3], s[2:3]
	s_xor_b64 exec, exec, s[2:3]
; %bb.291:
	v_ashrrev_i32_e32 v6, 31, v12
	v_add_u32_e32 v12, vcc, v12, v6
	v_xor_b32_e32 v6, v12, v6
; %bb.292:
	s_or_b64 exec, exec, s[2:3]
	v_add_u32_e32 v6, vcc, s7, v6
BB2_293:                                ; %Flow535
	v_add_u32_e32 v7, vcc, 2, v7
	v_and_b32_e32 v13, 1, v7
	v_mul_lo_u32 v13, v13, s20
	v_lshrrev_b32_e32 v12, 31, v7
	v_add_u32_e32 v7, vcc, v7, v12
	v_ashrrev_i32_e32 v7, 1, v7
	v_add_u32_e32 v13, vcc, v13, v7
	v_mov_b32_e32 v7, s12
BB2_294:
	s_or_b64 exec, exec, s[4:5]
	v_cvt_f32_u32_e32 v12, s22
	v_cvt_f32_u32_e32 v17, v0
	v_lshlrev_b64 v[15:16], 1, v[0:1]
	s_add_i32 s2, s15, -1
	v_rcp_iflag_f32_e32 v18, v12
	v_mov_b32_e32 v14, 0
	v_mov_b32_e32 v22, v14
	v_mul_f32_e32 v16, v17, v18
	v_trunc_f32_e32 v16, v16
	v_cvt_u32_f32_e32 v18, v16
	v_mad_f32 v16, -v16, v12, v17
	v_cmp_ge_f32_e64 vcc, |v16|, v12
	v_addc_u32_e32 v12, vcc, 0, v18, vcc
	v_and_b32_e32 v12, 0x3fffffff, v12
	v_mul_lo_u32 v16, v12, s2
	s_mul_i32 s2, s15, s6
	v_mov_b32_e32 v12, v14
	v_sub_u32_e32 v15, vcc, v15, v16
	v_add_u32_e32 v16, vcc, s2, v15
	v_cmp_gt_i32_e32 vcc, s12, v16
	s_and_saveexec_b64 s[2:3], vcc
	s_cbranch_execz BB2_300
; %bb.295:
	v_lshrrev_b32_e32 v12, 31, v16
	s_lshr_b32 s4, s12, 31
	s_lshr_b32 s5, s13, 31
	s_add_i32 s5, s13, s5
	v_add_u32_e32 v12, vcc, v16, v12
	s_add_i32 s4, s12, s4
	v_ashrrev_i32_e32 v14, 1, v12
	v_and_b32_e32 v12, 1, v16
	s_ashr_i32 s6, s4, 1
	s_and_b32 s4, s12, 1
	s_ashr_i32 s5, s5, 1
	s_and_b32 s7, s13, 1
	s_add_i32 s5, s5, s7
	v_cmp_eq_u32_e32 vcc, 1, v12
	s_add_i32 s4, s6, s4
                                        ; implicit-def: $vgpr16
                                        ; implicit-def: $sgpr7
	s_and_saveexec_b64 s[22:23], vcc
	s_xor_b64 s[22:23], exec, s[22:23]
; %bb.296:
	s_mul_i32 s7, s5, s4
	v_add_u32_e32 v16, vcc, s7, v14
	s_mul_i32 s7, s13, s12
	s_lshr_b32 s13, s7, 31
	s_add_i32 s7, s7, s13
	s_ashr_i32 s7, s7, 1
; %bb.297:                              ; %Flow529
	s_or_saveexec_b64 s[22:23], s[22:23]
	v_mov_b32_e32 v12, s7
	v_mov_b32_e32 v17, s6
	s_xor_b64 exec, exec, s[22:23]
; %bb.298:
	s_mul_i32 s5, s5, s12
	v_mov_b32_e32 v12, s5
	v_mov_b32_e32 v17, s4
	v_mov_b32_e32 v16, v14
; %bb.299:
	s_or_b64 exec, exec, s[22:23]
	s_lshr_b32 s4, s21, 31
	s_add_i32 s4, s21, s4
	s_and_b32 s5, s21, 1
	v_mul_lo_u32 v18, v12, s5
	s_ashr_i32 s4, s4, 1
	v_sub_u32_e32 v14, vcc, v17, v12
	v_mul_lo_u32 v17, v17, s4
	v_add_u32_e32 v16, vcc, v18, v16
	v_add_u32_e32 v22, vcc, v16, v17
BB2_300:                                ; %Flow530
	s_or_b64 exec, exec, s[2:3]
	v_add_u32_e32 v15, vcc, 2, v15
	v_and_b32_e32 v18, 1, v15
	v_mul_lo_u32 v17, v18, s20
	v_lshrrev_b32_e32 v16, 31, v15
	v_add_u32_e32 v15, vcc, v15, v16
	v_ashrrev_i32_e32 v21, 1, v15
	v_add_u32_e32 v15, vcc, v17, v21
	s_cmp_lt_i32 s14, 1
	s_cbranch_scc1 BB2_385
; %bb.301:
	s_sub_i32 s2, s18, s19
	v_lshlrev_b32_e32 v8, 2, v8
	s_add_i32 s2, s2, 32
	v_add_u32_e32 v16, vcc, 64, v8
	v_lshlrev_b32_e32 v8, 2, v13
	v_mul_lo_u32 v13, s2, v18
	v_add_u32_e32 v17, vcc, 64, v8
	v_ashrrev_i32_e32 v8, 31, v7
	s_ashr_i32 s17, s12, 31
	v_add_u32_e32 v20, vcc, v21, v13
	v_lshlrev_b64 v[18:19], 2, v[7:8]
	v_lshlrev_b32_e32 v13, 2, v13
	v_lshlrev_b32_e32 v21, 2, v21
	v_add_u32_e32 v21, vcc, v13, v21
	s_lshl_b64 s[6:7], s[16:17], 2
	s_mov_b32 s13, 0
	v_mov_b32_e32 v13, 0
	s_mov_b32 s15, 0x7ffffffc
	s_mov_b32 s18, 0x4f800000
	s_movk_i32 s19, 0xff
	s_branch BB2_304
BB2_302:                                ; %Flow485
                                        ;   in Loop: Header=BB2_304 Depth=1
	v_mov_b32_e32 v22, v23
BB2_303:                                ; %.loopexit.i25
                                        ;   in Loop: Header=BB2_304 Depth=1
	s_add_i32 s13, s13, 1
	s_cmp_eq_u32 s13, s14
	s_waitcnt vmcnt(0) lgkmcnt(0)
	s_barrier
	s_waitcnt lgkmcnt(0)
	s_cbranch_scc1 BB2_385
BB2_304:                                ; =>This Loop Header: Depth=1
                                        ;     Child Loop BB2_308 Depth 2
                                        ;     Child Loop BB2_312 Depth 2
                                        ;     Child Loop BB2_315 Depth 2
                                        ;     Child Loop BB2_319 Depth 2
                                        ;     Child Loop BB2_322 Depth 2
                                        ;     Child Loop BB2_328 Depth 2
                                        ;     Child Loop BB2_332 Depth 2
                                        ;     Child Loop BB2_335 Depth 2
                                        ;     Child Loop BB2_339 Depth 2
                                        ;     Child Loop BB2_342 Depth 2
                                        ;     Child Loop BB2_346 Depth 2
                                        ;     Child Loop BB2_353 Depth 2
                                        ;     Child Loop BB2_362 Depth 2
                                        ;     Child Loop BB2_382 Depth 2
	s_waitcnt vmcnt(0) lgkmcnt(0)
	ds_write_b32 v16, v5
	v_mov_b32_e32 v5, 0
	ds_read_b32 v23, v5 offset:8792
	s_mov_b64 s[2:3], -1
	s_waitcnt lgkmcnt(0)
	v_lshlrev_b32_e32 v23, 2, v23
	v_add_u32_e32 v23, vcc, v16, v23
	ds_write_b32 v23, v4
	ds_read_b32 v4, v5 offset:8792
	s_waitcnt lgkmcnt(0)
	v_lshlrev_b32_e32 v4, 3, v4
	v_add_u32_e32 v4, vcc, v16, v4
	ds_write_b32 v4, v3
	ds_read_b32 v3, v5 offset:4
                                        ; implicit-def: $vgpr4
	s_waitcnt lgkmcnt(0)
	v_cmp_gt_i32_e32 vcc, 1, v3
	s_and_b64 vcc, exec, vcc
	s_cbranch_vccz BB2_306
; %bb.305:                              ;   in Loop: Header=BB2_304 Depth=1
	v_add_u32_e32 v4, vcc, 2, v3
	s_mov_b64 s[2:3], 0
BB2_306:                                ; %Flow525
                                        ;   in Loop: Header=BB2_304 Depth=1
	s_andn2_b64 vcc, exec, s[2:3]
	s_cbranch_vccnz BB2_309
; %bb.307:                              ; %.preheader23.i7.preheader
                                        ;   in Loop: Header=BB2_304 Depth=1
	v_ashrrev_i32_e32 v4, 31, v2
	v_mov_b32_e32 v5, s17
	v_add_u32_e32 v3, vcc, s16, v2
	v_addc_u32_e32 v4, vcc, v5, v4, vcc
	v_lshlrev_b64 v[3:4], 2, v[3:4]
	v_mov_b32_e32 v5, s9
	v_add_u32_e32 v23, vcc, s8, v3
	s_mov_b32 s2, 2
	v_addc_u32_e32 v24, vcc, v5, v4, vcc
BB2_308:                                ; %.preheader23.i7
                                        ;   Parent Loop BB2_304 Depth=1
                                        ; =>  This Inner Loop Header: Depth=2
	flat_load_dword v3, v[23:24]
	v_mov_b32_e32 v4, 0
	ds_read_b32 v25, v4 offset:8792
	v_add_u32_e32 v2, vcc, s12, v2
	s_add_i32 s2, s2, 1
	v_mov_b32_e32 v5, s7
	v_add_u32_e32 v23, vcc, s6, v23
	v_addc_u32_e32 v24, vcc, v24, v5, vcc
	s_waitcnt lgkmcnt(0)
	v_mul_lo_u32 v5, v25, s2
	v_lshlrev_b32_e32 v5, 2, v5
	v_add_u32_e32 v5, vcc, v16, v5
	s_waitcnt vmcnt(0)
	ds_write_b32 v5, v3
	ds_read_b32 v3, v4 offset:4
	s_waitcnt lgkmcnt(0)
	v_add_u32_e32 v4, vcc, 2, v3
	v_cmp_ge_i32_e32 vcc, s2, v4
	s_and_b64 vcc, exec, vcc
	s_cbranch_vccz BB2_308
BB2_309:                                ; %.loopexit24.i3
                                        ;   in Loop: Header=BB2_304 Depth=1
	v_mov_b32_e32 v23, 0
	ds_read_b32 v5, v23 offset:8792
	ds_read_b32 v24, v23 offset:40
	s_waitcnt lgkmcnt(1)
	v_mul_lo_u32 v3, v5, v3
	v_mul_lo_u32 v4, v5, v4
	v_lshlrev_b32_e32 v5, 2, v5
	v_lshlrev_b32_e32 v3, 2, v3
	v_lshlrev_b32_e32 v4, 2, v4
	v_add_u32_e32 v25, vcc, v16, v3
	v_add_u32_e32 v3, vcc, v16, v4
	v_add_u32_e32 v4, vcc, v25, v5
	ds_read_b32 v3, v3
	ds_read_b32 v5, v25
	ds_read_b32 v4, v4
	s_waitcnt lgkmcnt(3)
	v_cmp_gt_i32_e32 vcc, 3, v24
	s_and_b64 vcc, exec, vcc
	s_cbranch_vccnz BB2_323
; %bb.310:                              ;   in Loop: Header=BB2_304 Depth=1
	v_add_u32_e32 v24, vcc, -1, v24
	v_lshrrev_b32_e32 v25, 31, v24
	v_add_u32_e32 v24, vcc, v25, v24
	v_ashrrev_i32_e32 v24, 1, v24
	v_max_i32_e32 v24, 1, v24
	v_add_u32_e32 v25, vcc, -1, v24
	v_cmp_gt_u32_e32 vcc, 3, v25
	s_and_b64 vcc, exec, vcc
	s_cbranch_vccnz BB2_313
; %bb.311:                              ;   in Loop: Header=BB2_304 Depth=1
	v_and_b32_e32 v23, s15, v24
	s_mov_b32 s4, 0
	s_mov_b32 s5, 8
BB2_312:                                ;   Parent Loop BB2_304 Depth=1
                                        ; =>  This Inner Loop Header: Depth=2
	v_mov_b32_e32 v25, 0
	ds_read_b32 v26, v25 offset:44
	s_add_i32 s20, s5, -8
	s_add_i32 s21, s5, -6
	s_add_i32 s22, s5, -4
	s_add_i32 s23, s5, -2
	s_waitcnt lgkmcnt(0)
	v_mul_lo_u32 v27, s21, v26
	v_mul_lo_u32 v28, v26, s20
	v_lshlrev_b32_e32 v26, 2, v26
	s_add_i32 s4, s4, 4
	v_lshlrev_b32_e32 v27, 2, v27
	v_lshlrev_b32_e32 v28, 2, v28
	v_add_u32_e64 v27, s[2:3], v16, v27
	v_add_u32_e64 v28, s[2:3], v16, v28
	ds_read_b32 v27, v27
	ds_read_b32 v29, v28
	v_add_u32_e64 v26, s[2:3], v28, v26
	ds_read_b32 v28, v26
	v_cmp_ne_u32_e32 vcc, s4, v23
	s_and_b64 vcc, exec, vcc
	s_waitcnt lgkmcnt(1)
	v_add_u32_e64 v27, s[2:3], v29, v27
	v_lshrrev_b32_e32 v29, 31, v27
	v_add_u32_e64 v27, s[2:3], v27, v29
	v_ashrrev_i32_e32 v27, 1, v27
	s_waitcnt lgkmcnt(0)
	v_sub_u32_e64 v27, s[2:3], v28, v27
	ds_write_b32 v26, v27
	ds_read_b32 v26, v25 offset:44
	s_waitcnt lgkmcnt(0)
	v_mul_lo_u32 v27, v26, s22
	v_mul_lo_u32 v28, v26, s21
	v_lshlrev_b32_e32 v26, 2, v26
	v_lshlrev_b32_e32 v27, 2, v27
	v_lshlrev_b32_e32 v28, 2, v28
	v_add_u32_e64 v27, s[2:3], v16, v27
	v_add_u32_e64 v28, s[2:3], v16, v28
	ds_read_b32 v27, v27
	ds_read_b32 v29, v28
	v_add_u32_e64 v26, s[2:3], v28, v26
	ds_read_b32 v28, v26
	s_waitcnt lgkmcnt(1)
	v_add_u32_e64 v27, s[2:3], v29, v27
	v_lshrrev_b32_e32 v29, 31, v27
	v_add_u32_e64 v27, s[2:3], v27, v29
	v_ashrrev_i32_e32 v27, 1, v27
	s_waitcnt lgkmcnt(0)
	v_sub_u32_e64 v27, s[2:3], v28, v27
	ds_write_b32 v26, v27
	ds_read_b32 v26, v25 offset:44
	s_waitcnt lgkmcnt(0)
	v_mul_lo_u32 v27, v26, s23
	v_mul_lo_u32 v28, v26, s22
	v_lshlrev_b32_e32 v26, 2, v26
	v_lshlrev_b32_e32 v27, 2, v27
	v_lshlrev_b32_e32 v28, 2, v28
	v_add_u32_e64 v27, s[2:3], v16, v27
	v_add_u32_e64 v28, s[2:3], v16, v28
	ds_read_b32 v27, v27
	ds_read_b32 v29, v28
	v_add_u32_e64 v26, s[2:3], v28, v26
	ds_read_b32 v28, v26
	s_waitcnt lgkmcnt(1)
	v_add_u32_e64 v27, s[2:3], v29, v27
	v_lshrrev_b32_e32 v29, 31, v27
	v_add_u32_e64 v27, s[2:3], v27, v29
	v_ashrrev_i32_e32 v27, 1, v27
	s_waitcnt lgkmcnt(0)
	v_sub_u32_e64 v27, s[2:3], v28, v27
	ds_write_b32 v26, v27
	ds_read_b32 v25, v25 offset:44
	s_waitcnt lgkmcnt(0)
	v_mul_lo_u32 v26, v25, s5
	v_mul_lo_u32 v27, v25, s23
	v_lshlrev_b32_e32 v25, 2, v25
	s_add_i32 s5, s5, 8
	v_lshlrev_b32_e32 v26, 2, v26
	v_lshlrev_b32_e32 v27, 2, v27
	v_add_u32_e64 v26, s[2:3], v16, v26
	v_add_u32_e64 v27, s[2:3], v16, v27
	ds_read_b32 v26, v26
	ds_read_b32 v28, v27
	v_add_u32_e64 v25, s[2:3], v27, v25
	ds_read_b32 v27, v25
	s_waitcnt lgkmcnt(1)
	v_add_u32_e64 v26, s[2:3], v28, v26
	v_lshrrev_b32_e32 v28, 31, v26
	v_add_u32_e64 v26, s[2:3], v26, v28
	v_ashrrev_i32_e32 v26, 1, v26
	s_waitcnt lgkmcnt(0)
	v_sub_u32_e64 v26, s[2:3], v27, v26
	ds_write_b32 v25, v26
	s_cbranch_vccnz BB2_312
BB2_313:                                ; %Flow522
                                        ;   in Loop: Header=BB2_304 Depth=1
	v_and_b32_e32 v24, 3, v24
	v_cmp_eq_u32_e32 vcc, 0, v24
	s_and_b64 vcc, exec, vcc
	s_cbranch_vccnz BB2_316
; %bb.314:                              ; %.preheader20.i5.preheader
                                        ;   in Loop: Header=BB2_304 Depth=1
	v_lshlrev_b32_e32 v23, 1, v23
BB2_315:                                ; %.preheader20.i5
                                        ;   Parent Loop BB2_304 Depth=1
                                        ; =>  This Inner Loop Header: Depth=2
	v_mov_b32_e32 v25, 0
	ds_read_b32 v25, v25 offset:44
	v_add_u32_e32 v26, vcc, 2, v23
	v_add_u32_e32 v24, vcc, -1, v24
	v_cmp_ne_u32_e32 vcc, 0, v24
	s_and_b64 vcc, exec, vcc
	s_waitcnt lgkmcnt(0)
	v_mul_lo_u32 v27, v26, v25
	v_mul_lo_u32 v28, v25, v23
	v_mov_b32_e32 v23, v26
	v_lshlrev_b32_e32 v25, 2, v25
	v_lshlrev_b32_e32 v27, 2, v27
	v_lshlrev_b32_e32 v26, 2, v28
	v_add_u32_e64 v27, s[2:3], v16, v27
	v_add_u32_e64 v26, s[2:3], v16, v26
	ds_read_b32 v27, v27
	ds_read_b32 v28, v26
	v_add_u32_e64 v25, s[2:3], v26, v25
	ds_read_b32 v26, v25
	s_waitcnt lgkmcnt(1)
	v_add_u32_e64 v27, s[2:3], v28, v27
	v_lshrrev_b32_e32 v28, 31, v27
	v_add_u32_e64 v27, s[2:3], v27, v28
	v_ashrrev_i32_e32 v27, 1, v27
	s_waitcnt lgkmcnt(0)
	v_sub_u32_e64 v26, s[2:3], v26, v27
	ds_write_b32 v25, v26
	s_cbranch_vccnz BB2_315
BB2_316:                                ; %.loopexit21.i6
                                        ;   in Loop: Header=BB2_304 Depth=1
	v_mov_b32_e32 v23, 0
	ds_read_b32 v24, v23 offset:40
	s_waitcnt lgkmcnt(0)
	v_cmp_gt_i32_e32 vcc, 4, v24
	s_and_b64 vcc, exec, vcc
	s_cbranch_vccnz BB2_323
; %bb.317:                              ;   in Loop: Header=BB2_304 Depth=1
	v_lshrrev_b32_e32 v24, 1, v24
	v_max_u32_e32 v25, 2, v24
	v_add_u32_e32 v24, vcc, -1, v25
	v_add_u32_e32 v25, vcc, -2, v25
	v_cmp_gt_u32_e32 vcc, 3, v25
	s_and_b64 vcc, exec, vcc
	s_cbranch_vccnz BB2_320
; %bb.318:                              ;   in Loop: Header=BB2_304 Depth=1
	v_and_b32_e32 v23, -4, v24
	s_mov_b32 s4, 0
	s_mov_b32 s5, 9
BB2_319:                                ;   Parent Loop BB2_304 Depth=1
                                        ; =>  This Inner Loop Header: Depth=2
	v_mov_b32_e32 v25, 0
	ds_read_b32 v26, v25 offset:44
	s_add_i32 s2, s5, -6
	s_add_i32 s3, s5, -4
	s_add_i32 s20, s5, -2
	s_add_i32 s4, s4, 4
	s_waitcnt lgkmcnt(0)
	v_mul_lo_u32 v27, v26, s2
	v_lshlrev_b32_e32 v29, 1, v26
	v_lshlrev_b32_e32 v26, 2, v26
	v_lshlrev_b32_e32 v28, 2, v27
	v_subrev_u32_e32 v27, vcc, v29, v27
	v_lshlrev_b32_e32 v27, 2, v27
	v_add_u32_e32 v28, vcc, v16, v28
	v_add_u32_e32 v27, vcc, v16, v27
	ds_read_b32 v28, v28
	ds_read_b32 v29, v27
	v_add_u32_e32 v26, vcc, v27, v26
	ds_read_b32 v27, v26
	s_waitcnt lgkmcnt(1)
	v_add_u32_e32 v28, vcc, v28, v29
	v_add_u32_e32 v28, vcc, 2, v28
	v_ashrrev_i32_e32 v29, 31, v28
	v_lshrrev_b32_e32 v29, 30, v29
	v_add_u32_e32 v28, vcc, v28, v29
	v_ashrrev_i32_e32 v28, 2, v28
	s_waitcnt lgkmcnt(0)
	v_add_u32_e32 v27, vcc, v28, v27
	ds_write_b32 v26, v27
	ds_read_b32 v26, v25 offset:44
	s_waitcnt lgkmcnt(0)
	v_mul_lo_u32 v29, v26, s2
	v_mul_lo_u32 v27, v26, s3
	v_lshlrev_b32_e32 v29, 2, v29
	v_lshlrev_b32_e32 v28, 2, v27
	v_add_u32_e32 v28, vcc, v16, v28
	v_add_u32_e32 v29, vcc, v16, v29
	ds_read_b32 v28, v28
	ds_read_b32 v29, v29
	v_subrev_u32_e32 v26, vcc, v26, v27
	v_lshlrev_b32_e32 v26, 2, v26
	v_add_u32_e32 v26, vcc, v16, v26
	ds_read_b32 v27, v26
	s_waitcnt lgkmcnt(1)
	v_add_u32_e32 v28, vcc, v28, v29
	v_add_u32_e32 v28, vcc, 2, v28
	v_ashrrev_i32_e32 v29, 31, v28
	v_lshrrev_b32_e32 v29, 30, v29
	v_add_u32_e32 v28, vcc, v28, v29
	v_ashrrev_i32_e32 v28, 2, v28
	s_waitcnt lgkmcnt(0)
	v_add_u32_e32 v27, vcc, v28, v27
	ds_write_b32 v26, v27
	ds_read_b32 v26, v25 offset:44
	s_waitcnt lgkmcnt(0)
	v_mul_lo_u32 v27, v26, s20
	v_lshlrev_b32_e32 v29, 1, v26
	v_lshlrev_b32_e32 v26, 2, v26
	v_lshlrev_b32_e32 v28, 2, v27
	v_subrev_u32_e32 v27, vcc, v29, v27
	v_lshlrev_b32_e32 v27, 2, v27
	v_add_u32_e32 v28, vcc, v16, v28
	v_add_u32_e32 v27, vcc, v16, v27
	ds_read_b32 v28, v28
	ds_read_b32 v29, v27
	v_add_u32_e32 v26, vcc, v27, v26
	ds_read_b32 v27, v26
	s_waitcnt lgkmcnt(1)
	v_add_u32_e32 v28, vcc, v28, v29
	v_add_u32_e32 v28, vcc, 2, v28
	v_ashrrev_i32_e32 v29, 31, v28
	v_lshrrev_b32_e32 v29, 30, v29
	v_add_u32_e32 v28, vcc, v28, v29
	v_ashrrev_i32_e32 v28, 2, v28
	s_waitcnt lgkmcnt(0)
	v_add_u32_e32 v27, vcc, v28, v27
	ds_write_b32 v26, v27
	ds_read_b32 v25, v25 offset:44
	s_waitcnt lgkmcnt(0)
	v_mul_lo_u32 v26, v25, s5
	v_mul_lo_u32 v28, v25, s20
	s_add_i32 s5, s5, 8
	v_lshlrev_b32_e32 v27, 2, v26
	v_lshlrev_b32_e32 v28, 2, v28
	v_add_u32_e32 v27, vcc, v16, v27
	v_add_u32_e32 v28, vcc, v16, v28
	ds_read_b32 v27, v27
	ds_read_b32 v28, v28
	v_subrev_u32_e32 v25, vcc, v25, v26
	v_lshlrev_b32_e32 v25, 2, v25
	v_add_u32_e64 v25, s[2:3], v16, v25
	ds_read_b32 v26, v25
	s_waitcnt lgkmcnt(1)
	v_add_u32_e64 v27, s[2:3], v27, v28
	v_add_u32_e64 v27, s[2:3], 2, v27
	v_ashrrev_i32_e32 v28, 31, v27
	v_lshrrev_b32_e32 v28, 30, v28
	v_add_u32_e64 v27, s[2:3], v27, v28
	v_cmp_eq_u32_e32 vcc, s4, v23
	v_ashrrev_i32_e32 v27, 2, v27
	s_and_b64 vcc, exec, vcc
	s_waitcnt lgkmcnt(0)
	v_add_u32_e64 v26, s[2:3], v27, v26
	ds_write_b32 v25, v26
	s_cbranch_vccz BB2_319
BB2_320:                                ; %Flow517
                                        ;   in Loop: Header=BB2_304 Depth=1
	v_and_b32_e32 v24, 3, v24
	v_cmp_eq_u32_e32 vcc, 0, v24
	s_and_b64 vcc, exec, vcc
	s_cbranch_vccnz BB2_323
; %bb.321:                              ; %.preheader17.i9.preheader
                                        ;   in Loop: Header=BB2_304 Depth=1
	v_lshlrev_b32_e32 v23, 1, v23
	v_add_u32_e32 v23, vcc, 3, v23
BB2_322:                                ; %.preheader17.i9
                                        ;   Parent Loop BB2_304 Depth=1
                                        ; =>  This Inner Loop Header: Depth=2
	v_mov_b32_e32 v25, 0
	ds_read_b32 v25, v25 offset:44
	v_add_u32_e32 v24, vcc, -1, v24
	s_waitcnt lgkmcnt(0)
	v_mul_lo_u32 v26, v25, v23
	v_lshlrev_b32_e32 v28, 1, v25
	v_lshlrev_b32_e32 v25, 2, v25
	v_add_u32_e32 v23, vcc, 2, v23
	v_lshlrev_b32_e32 v27, 2, v26
	v_subrev_u32_e32 v26, vcc, v28, v26
	v_lshlrev_b32_e32 v26, 2, v26
	v_add_u32_e32 v27, vcc, v16, v27
	v_add_u32_e64 v26, s[2:3], v16, v26
	ds_read_b32 v27, v27
	ds_read_b32 v28, v26
	v_add_u32_e64 v25, s[2:3], v26, v25
	ds_read_b32 v26, v25
	v_cmp_ne_u32_e32 vcc, 0, v24
	s_and_b64 vcc, exec, vcc
	s_waitcnt lgkmcnt(1)
	v_add_u32_e64 v27, s[2:3], v27, v28
	v_add_u32_e64 v27, s[2:3], 2, v27
	v_ashrrev_i32_e32 v28, 31, v27
	v_lshrrev_b32_e32 v28, 30, v28
	v_add_u32_e64 v27, s[2:3], v27, v28
	v_ashrrev_i32_e32 v27, 2, v27
	s_waitcnt lgkmcnt(0)
	v_add_u32_e64 v26, s[2:3], v27, v26
	ds_write_b32 v25, v26
	s_cbranch_vccnz BB2_322
BB2_323:                                ; %.loopexit18.i10
                                        ;   in Loop: Header=BB2_304 Depth=1
	s_and_saveexec_b64 s[4:5], s[0:1]
	s_cbranch_execz BB2_343
; %bb.324:                              ;   in Loop: Header=BB2_304 Depth=1
	ds_write_b32 v17, v11
	v_mov_b32_e32 v11, 0
	ds_read_b32 v23, v11 offset:8792
	s_waitcnt lgkmcnt(0)
	v_lshlrev_b32_e32 v23, 2, v23
	v_add_u32_e32 v23, vcc, v17, v23
	ds_write_b32 v23, v10
	ds_read_b32 v10, v11 offset:8792
	s_waitcnt lgkmcnt(0)
	v_lshlrev_b32_e32 v10, 3, v10
	v_add_u32_e32 v10, vcc, v17, v10
	ds_write_b32 v10, v9
	ds_read_b32 v9, v11 offset:4
	s_waitcnt lgkmcnt(0)
	v_cmp_gt_i32_e32 vcc, 1, v9
	s_and_b64 vcc, exec, vcc
	s_cbranch_vccz BB2_326
; %bb.325:                              ;   in Loop: Header=BB2_304 Depth=1
	v_add_u32_e32 v10, vcc, 2, v9
	s_mov_b64 s[2:3], 0
	s_andn2_b64 vcc, exec, s[2:3]
	s_cbranch_vccz BB2_327
	s_branch BB2_329
BB2_326:                                ;   in Loop: Header=BB2_304 Depth=1
	s_mov_b64 s[2:3], -1
                                        ; implicit-def: $vgpr10
	s_andn2_b64 vcc, exec, s[2:3]
	s_cbranch_vccnz BB2_329
BB2_327:                                ; %.preheader15.i15.preheader
                                        ;   in Loop: Header=BB2_304 Depth=1
	v_ashrrev_i32_e32 v10, 31, v6
	v_add_u32_e32 v9, vcc, v7, v6
	v_addc_u32_e32 v10, vcc, v8, v10, vcc
	v_lshlrev_b64 v[9:10], 2, v[9:10]
	v_mov_b32_e32 v11, s9
	v_add_u32_e32 v23, vcc, s8, v9
	s_mov_b32 s2, 2
	v_addc_u32_e32 v24, vcc, v11, v10, vcc
BB2_328:                                ; %.preheader15.i15
                                        ;   Parent Loop BB2_304 Depth=1
                                        ; =>  This Inner Loop Header: Depth=2
	flat_load_dword v9, v[23:24]
	v_mov_b32_e32 v10, 0
	ds_read_b32 v11, v10 offset:8792
	s_add_i32 s2, s2, 1
	v_add_u32_e32 v6, vcc, v6, v7
	v_add_u32_e32 v23, vcc, v23, v18
	v_addc_u32_e32 v24, vcc, v24, v19, vcc
	s_waitcnt lgkmcnt(0)
	v_mul_lo_u32 v11, v11, s2
	v_lshlrev_b32_e32 v11, 2, v11
	v_add_u32_e32 v11, vcc, v17, v11
	s_waitcnt vmcnt(0)
	ds_write_b32 v11, v9
	ds_read_b32 v9, v10 offset:4
	s_waitcnt lgkmcnt(0)
	v_add_u32_e32 v10, vcc, 2, v9
	v_cmp_ge_i32_e32 vcc, s2, v10
	s_and_b64 vcc, exec, vcc
	s_cbranch_vccz BB2_328
BB2_329:                                ; %.loopexit16.i11
                                        ;   in Loop: Header=BB2_304 Depth=1
	v_mov_b32_e32 v23, 0
	ds_read_b32 v11, v23 offset:8792
	ds_read_b32 v24, v23 offset:40
	s_waitcnt lgkmcnt(1)
	v_mul_lo_u32 v9, v11, v9
	v_mul_lo_u32 v10, v11, v10
	v_lshlrev_b32_e32 v11, 2, v11
	v_lshlrev_b32_e32 v9, 2, v9
	v_lshlrev_b32_e32 v10, 2, v10
	v_add_u32_e32 v25, vcc, v17, v9
	v_add_u32_e32 v9, vcc, v17, v10
	v_add_u32_e32 v10, vcc, v25, v11
	ds_read_b32 v9, v9
	ds_read_b32 v11, v25
	ds_read_b32 v10, v10
	s_waitcnt lgkmcnt(3)
	v_cmp_gt_i32_e32 vcc, 3, v24
	s_and_b64 vcc, exec, vcc
	s_cbranch_vccnz BB2_343
; %bb.330:                              ;   in Loop: Header=BB2_304 Depth=1
	v_add_u32_e32 v24, vcc, -1, v24
	v_lshrrev_b32_e32 v25, 31, v24
	v_add_u32_e32 v24, vcc, v25, v24
	v_ashrrev_i32_e32 v24, 1, v24
	v_max_i32_e32 v24, 1, v24
	v_add_u32_e32 v25, vcc, -1, v24
	v_cmp_gt_u32_e32 vcc, 3, v25
	s_and_b64 vcc, exec, vcc
	s_cbranch_vccnz BB2_333
; %bb.331:                              ;   in Loop: Header=BB2_304 Depth=1
	v_and_b32_e32 v23, s15, v24
	s_mov_b32 s20, 0
	s_mov_b32 s21, 8
BB2_332:                                ;   Parent Loop BB2_304 Depth=1
                                        ; =>  This Inner Loop Header: Depth=2
	v_mov_b32_e32 v25, 0
	ds_read_b32 v26, v25 offset:44
	s_add_i32 s22, s21, -8
	s_add_i32 s23, s21, -6
	s_add_i32 s24, s21, -4
	s_add_i32 s25, s21, -2
	s_waitcnt lgkmcnt(0)
	v_mul_lo_u32 v27, s23, v26
	v_mul_lo_u32 v28, v26, s22
	v_lshlrev_b32_e32 v26, 2, v26
	s_add_i32 s20, s20, 4
	v_lshlrev_b32_e32 v27, 2, v27
	v_lshlrev_b32_e32 v28, 2, v28
	v_add_u32_e64 v27, s[2:3], v17, v27
	v_add_u32_e64 v28, s[2:3], v17, v28
	ds_read_b32 v27, v27
	ds_read_b32 v29, v28
	v_add_u32_e64 v26, s[2:3], v28, v26
	ds_read_b32 v28, v26
	v_cmp_ne_u32_e32 vcc, s20, v23
	s_and_b64 vcc, exec, vcc
	s_waitcnt lgkmcnt(1)
	v_add_u32_e64 v27, s[2:3], v29, v27
	v_lshrrev_b32_e32 v29, 31, v27
	v_add_u32_e64 v27, s[2:3], v27, v29
	v_ashrrev_i32_e32 v27, 1, v27
	s_waitcnt lgkmcnt(0)
	v_sub_u32_e64 v27, s[2:3], v28, v27
	ds_write_b32 v26, v27
	ds_read_b32 v26, v25 offset:44
	s_waitcnt lgkmcnt(0)
	v_mul_lo_u32 v27, v26, s24
	v_mul_lo_u32 v28, v26, s23
	v_lshlrev_b32_e32 v26, 2, v26
	v_lshlrev_b32_e32 v27, 2, v27
	v_lshlrev_b32_e32 v28, 2, v28
	v_add_u32_e64 v27, s[2:3], v17, v27
	v_add_u32_e64 v28, s[2:3], v17, v28
	ds_read_b32 v27, v27
	ds_read_b32 v29, v28
	v_add_u32_e64 v26, s[2:3], v28, v26
	ds_read_b32 v28, v26
	s_waitcnt lgkmcnt(1)
	v_add_u32_e64 v27, s[2:3], v29, v27
	v_lshrrev_b32_e32 v29, 31, v27
	v_add_u32_e64 v27, s[2:3], v27, v29
	v_ashrrev_i32_e32 v27, 1, v27
	s_waitcnt lgkmcnt(0)
	v_sub_u32_e64 v27, s[2:3], v28, v27
	ds_write_b32 v26, v27
	ds_read_b32 v26, v25 offset:44
	s_waitcnt lgkmcnt(0)
	v_mul_lo_u32 v27, v26, s25
	v_mul_lo_u32 v28, v26, s24
	v_lshlrev_b32_e32 v26, 2, v26
	v_lshlrev_b32_e32 v27, 2, v27
	v_lshlrev_b32_e32 v28, 2, v28
	v_add_u32_e64 v27, s[2:3], v17, v27
	v_add_u32_e64 v28, s[2:3], v17, v28
	ds_read_b32 v27, v27
	ds_read_b32 v29, v28
	v_add_u32_e64 v26, s[2:3], v28, v26
	ds_read_b32 v28, v26
	s_waitcnt lgkmcnt(1)
	v_add_u32_e64 v27, s[2:3], v29, v27
	v_lshrrev_b32_e32 v29, 31, v27
	v_add_u32_e64 v27, s[2:3], v27, v29
	v_ashrrev_i32_e32 v27, 1, v27
	s_waitcnt lgkmcnt(0)
	v_sub_u32_e64 v27, s[2:3], v28, v27
	ds_write_b32 v26, v27
	ds_read_b32 v25, v25 offset:44
	s_waitcnt lgkmcnt(0)
	v_mul_lo_u32 v26, v25, s21
	v_mul_lo_u32 v27, v25, s25
	v_lshlrev_b32_e32 v25, 2, v25
	s_add_i32 s21, s21, 8
	v_lshlrev_b32_e32 v26, 2, v26
	v_lshlrev_b32_e32 v27, 2, v27
	v_add_u32_e64 v26, s[2:3], v17, v26
	v_add_u32_e64 v27, s[2:3], v17, v27
	ds_read_b32 v26, v26
	ds_read_b32 v28, v27
	v_add_u32_e64 v25, s[2:3], v27, v25
	ds_read_b32 v27, v25
	s_waitcnt lgkmcnt(1)
	v_add_u32_e64 v26, s[2:3], v28, v26
	v_lshrrev_b32_e32 v28, 31, v26
	v_add_u32_e64 v26, s[2:3], v26, v28
	v_ashrrev_i32_e32 v26, 1, v26
	s_waitcnt lgkmcnt(0)
	v_sub_u32_e64 v26, s[2:3], v27, v26
	ds_write_b32 v25, v26
	s_cbranch_vccnz BB2_332
BB2_333:                                ; %Flow508
                                        ;   in Loop: Header=BB2_304 Depth=1
	v_and_b32_e32 v24, 3, v24
	v_cmp_eq_u32_e32 vcc, 0, v24
	s_and_b64 vcc, exec, vcc
	s_cbranch_vccnz BB2_336
; %bb.334:                              ; %.preheader12.i13.preheader
                                        ;   in Loop: Header=BB2_304 Depth=1
	v_lshlrev_b32_e32 v23, 1, v23
BB2_335:                                ; %.preheader12.i13
                                        ;   Parent Loop BB2_304 Depth=1
                                        ; =>  This Inner Loop Header: Depth=2
	v_mov_b32_e32 v25, 0
	ds_read_b32 v25, v25 offset:44
	v_add_u32_e32 v26, vcc, 2, v23
	v_add_u32_e32 v24, vcc, -1, v24
	v_cmp_ne_u32_e32 vcc, 0, v24
	s_and_b64 vcc, exec, vcc
	s_waitcnt lgkmcnt(0)
	v_mul_lo_u32 v27, v26, v25
	v_mul_lo_u32 v28, v25, v23
	v_mov_b32_e32 v23, v26
	v_lshlrev_b32_e32 v25, 2, v25
	v_lshlrev_b32_e32 v27, 2, v27
	v_lshlrev_b32_e32 v26, 2, v28
	v_add_u32_e64 v27, s[2:3], v17, v27
	v_add_u32_e64 v26, s[2:3], v17, v26
	ds_read_b32 v27, v27
	ds_read_b32 v28, v26
	v_add_u32_e64 v25, s[2:3], v26, v25
	ds_read_b32 v26, v25
	s_waitcnt lgkmcnt(1)
	v_add_u32_e64 v27, s[2:3], v28, v27
	v_lshrrev_b32_e32 v28, 31, v27
	v_add_u32_e64 v27, s[2:3], v27, v28
	v_ashrrev_i32_e32 v27, 1, v27
	s_waitcnt lgkmcnt(0)
	v_sub_u32_e64 v26, s[2:3], v26, v27
	ds_write_b32 v25, v26
	s_cbranch_vccnz BB2_335
BB2_336:                                ; %.loopexit13.i14
                                        ;   in Loop: Header=BB2_304 Depth=1
	v_mov_b32_e32 v23, 0
	ds_read_b32 v24, v23 offset:40
	s_waitcnt lgkmcnt(0)
	v_cmp_gt_i32_e32 vcc, 4, v24
	s_and_b64 vcc, exec, vcc
	s_cbranch_vccnz BB2_343
; %bb.337:                              ;   in Loop: Header=BB2_304 Depth=1
	v_lshrrev_b32_e32 v24, 1, v24
	v_max_u32_e32 v25, 2, v24
	v_add_u32_e32 v24, vcc, -1, v25
	v_add_u32_e32 v25, vcc, -2, v25
	v_cmp_gt_u32_e32 vcc, 3, v25
	s_and_b64 vcc, exec, vcc
	s_cbranch_vccnz BB2_340
; %bb.338:                              ;   in Loop: Header=BB2_304 Depth=1
	v_and_b32_e32 v23, -4, v24
	s_mov_b32 s20, 0
	s_mov_b32 s21, 9
BB2_339:                                ;   Parent Loop BB2_304 Depth=1
                                        ; =>  This Inner Loop Header: Depth=2
	v_mov_b32_e32 v25, 0
	ds_read_b32 v26, v25 offset:44
	s_add_i32 s2, s21, -6
	s_add_i32 s3, s21, -4
	s_add_i32 s22, s21, -2
	s_add_i32 s20, s20, 4
	s_waitcnt lgkmcnt(0)
	v_mul_lo_u32 v27, v26, s2
	v_lshlrev_b32_e32 v29, 1, v26
	v_lshlrev_b32_e32 v26, 2, v26
	v_lshlrev_b32_e32 v28, 2, v27
	v_subrev_u32_e32 v27, vcc, v29, v27
	v_lshlrev_b32_e32 v27, 2, v27
	v_add_u32_e32 v28, vcc, v17, v28
	v_add_u32_e32 v27, vcc, v17, v27
	ds_read_b32 v28, v28
	ds_read_b32 v29, v27
	v_add_u32_e32 v26, vcc, v27, v26
	ds_read_b32 v27, v26
	s_waitcnt lgkmcnt(1)
	v_add_u32_e32 v28, vcc, v28, v29
	v_add_u32_e32 v28, vcc, 2, v28
	v_ashrrev_i32_e32 v29, 31, v28
	v_lshrrev_b32_e32 v29, 30, v29
	v_add_u32_e32 v28, vcc, v28, v29
	v_ashrrev_i32_e32 v28, 2, v28
	s_waitcnt lgkmcnt(0)
	v_add_u32_e32 v27, vcc, v28, v27
	ds_write_b32 v26, v27
	ds_read_b32 v26, v25 offset:44
	s_waitcnt lgkmcnt(0)
	v_mul_lo_u32 v29, v26, s2
	v_mul_lo_u32 v27, v26, s3
	v_lshlrev_b32_e32 v29, 2, v29
	v_lshlrev_b32_e32 v28, 2, v27
	v_add_u32_e32 v28, vcc, v17, v28
	v_add_u32_e32 v29, vcc, v17, v29
	ds_read_b32 v28, v28
	ds_read_b32 v29, v29
	v_subrev_u32_e32 v26, vcc, v26, v27
	v_lshlrev_b32_e32 v26, 2, v26
	v_add_u32_e32 v26, vcc, v17, v26
	ds_read_b32 v27, v26
	s_waitcnt lgkmcnt(1)
	v_add_u32_e32 v28, vcc, v28, v29
	v_add_u32_e32 v28, vcc, 2, v28
	v_ashrrev_i32_e32 v29, 31, v28
	v_lshrrev_b32_e32 v29, 30, v29
	v_add_u32_e32 v28, vcc, v28, v29
	v_ashrrev_i32_e32 v28, 2, v28
	s_waitcnt lgkmcnt(0)
	v_add_u32_e32 v27, vcc, v28, v27
	ds_write_b32 v26, v27
	ds_read_b32 v26, v25 offset:44
	s_waitcnt lgkmcnt(0)
	v_mul_lo_u32 v27, v26, s22
	v_lshlrev_b32_e32 v29, 1, v26
	v_lshlrev_b32_e32 v26, 2, v26
	v_lshlrev_b32_e32 v28, 2, v27
	v_subrev_u32_e32 v27, vcc, v29, v27
	v_lshlrev_b32_e32 v27, 2, v27
	v_add_u32_e32 v28, vcc, v17, v28
	v_add_u32_e32 v27, vcc, v17, v27
	ds_read_b32 v28, v28
	ds_read_b32 v29, v27
	v_add_u32_e32 v26, vcc, v27, v26
	ds_read_b32 v27, v26
	s_waitcnt lgkmcnt(1)
	v_add_u32_e32 v28, vcc, v28, v29
	v_add_u32_e32 v28, vcc, 2, v28
	v_ashrrev_i32_e32 v29, 31, v28
	v_lshrrev_b32_e32 v29, 30, v29
	v_add_u32_e32 v28, vcc, v28, v29
	v_ashrrev_i32_e32 v28, 2, v28
	s_waitcnt lgkmcnt(0)
	v_add_u32_e32 v27, vcc, v28, v27
	ds_write_b32 v26, v27
	ds_read_b32 v25, v25 offset:44
	s_waitcnt lgkmcnt(0)
	v_mul_lo_u32 v26, v25, s21
	v_mul_lo_u32 v28, v25, s22
	s_add_i32 s21, s21, 8
	v_lshlrev_b32_e32 v27, 2, v26
	v_lshlrev_b32_e32 v28, 2, v28
	v_add_u32_e32 v27, vcc, v17, v27
	v_add_u32_e32 v28, vcc, v17, v28
	ds_read_b32 v27, v27
	ds_read_b32 v28, v28
	v_subrev_u32_e32 v25, vcc, v25, v26
	v_lshlrev_b32_e32 v25, 2, v25
	v_add_u32_e64 v25, s[2:3], v17, v25
	ds_read_b32 v26, v25
	s_waitcnt lgkmcnt(1)
	v_add_u32_e64 v27, s[2:3], v27, v28
	v_add_u32_e64 v27, s[2:3], 2, v27
	v_ashrrev_i32_e32 v28, 31, v27
	v_lshrrev_b32_e32 v28, 30, v28
	v_add_u32_e64 v27, s[2:3], v27, v28
	v_cmp_eq_u32_e32 vcc, s20, v23
	v_ashrrev_i32_e32 v27, 2, v27
	s_and_b64 vcc, exec, vcc
	s_waitcnt lgkmcnt(0)
	v_add_u32_e64 v26, s[2:3], v27, v26
	ds_write_b32 v25, v26
	s_cbranch_vccz BB2_339
BB2_340:                                ; %Flow503
                                        ;   in Loop: Header=BB2_304 Depth=1
	v_and_b32_e32 v24, 3, v24
	v_cmp_eq_u32_e32 vcc, 0, v24
	s_and_b64 vcc, exec, vcc
	s_cbranch_vccnz BB2_343
; %bb.341:                              ; %.preheader9.i17.preheader
                                        ;   in Loop: Header=BB2_304 Depth=1
	v_lshlrev_b32_e32 v23, 1, v23
	v_add_u32_e32 v23, vcc, 3, v23
BB2_342:                                ; %.preheader9.i17
                                        ;   Parent Loop BB2_304 Depth=1
                                        ; =>  This Inner Loop Header: Depth=2
	v_mov_b32_e32 v25, 0
	ds_read_b32 v25, v25 offset:44
	v_add_u32_e32 v24, vcc, -1, v24
	s_waitcnt lgkmcnt(0)
	v_mul_lo_u32 v26, v25, v23
	v_lshlrev_b32_e32 v28, 1, v25
	v_lshlrev_b32_e32 v25, 2, v25
	v_add_u32_e32 v23, vcc, 2, v23
	v_lshlrev_b32_e32 v27, 2, v26
	v_subrev_u32_e32 v26, vcc, v28, v26
	v_lshlrev_b32_e32 v26, 2, v26
	v_add_u32_e32 v27, vcc, v17, v27
	v_add_u32_e64 v26, s[2:3], v17, v26
	ds_read_b32 v27, v27
	ds_read_b32 v28, v26
	v_add_u32_e64 v25, s[2:3], v26, v25
	ds_read_b32 v26, v25
	v_cmp_ne_u32_e32 vcc, 0, v24
	s_and_b64 vcc, exec, vcc
	s_waitcnt lgkmcnt(1)
	v_add_u32_e64 v27, s[2:3], v27, v28
	v_add_u32_e64 v27, s[2:3], 2, v27
	v_ashrrev_i32_e32 v28, 31, v27
	v_lshrrev_b32_e32 v28, 30, v28
	v_add_u32_e64 v27, s[2:3], v27, v28
	v_ashrrev_i32_e32 v27, 2, v27
	s_waitcnt lgkmcnt(0)
	v_add_u32_e64 v26, s[2:3], v27, v26
	ds_write_b32 v25, v26
	s_cbranch_vccnz BB2_342
BB2_343:                                ; %Flow513
                                        ;   in Loop: Header=BB2_304 Depth=1
	s_or_b64 exec, exec, s[4:5]
	v_mov_b32_e32 v25, 0
	s_waitcnt lgkmcnt(0)
	s_barrier
	s_waitcnt lgkmcnt(0)
	ds_read2_b32 v[23:24], v25 offset0:1 offset1:9
	s_waitcnt lgkmcnt(0)
	v_ashrrev_i32_e32 v28, 31, v24
	v_add_u32_e32 v26, vcc, v28, v24
	v_xor_b32_e32 v29, v26, v28
	v_cvt_f32_u32_e32 v26, v29
	v_rcp_iflag_f32_e32 v30, v26
	ds_read2_b32 v[26:27], v25 offset0:11 offset1:15
	s_waitcnt lgkmcnt(0)
	v_mul_lo_u32 v25, v26, v23
	v_mul_f32_e32 v23, s18, v30
	v_cvt_u32_f32_e32 v30, v23
	v_lshlrev_b32_e32 v23, 1, v26
	v_add_u32_e32 v26, vcc, -1, v25
	v_ashrrev_i32_e32 v32, 31, v26
	v_mul_lo_u32 v25, v30, v29
	v_mul_hi_u32 v31, v30, v29
	v_xor_b32_e32 v28, v32, v28
	v_sub_u32_e32 v33, vcc, 0, v25
	v_cmp_eq_u32_e64 s[2:3], 0, v31
	v_cndmask_b32_e64 v25, v25, v33, s[2:3]
	v_mul_hi_u32 v25, v25, v30
	v_add_u32_e32 v31, vcc, v32, v26
	v_xor_b32_e32 v31, v31, v32
	v_add_u32_e32 v33, vcc, v25, v30
	v_subrev_u32_e32 v25, vcc, v25, v30
	v_cndmask_b32_e64 v25, v25, v33, s[2:3]
	v_mul_hi_u32 v30, v25, v31
	v_add_u32_e32 v25, vcc, v27, v23
	v_mul_lo_u32 v27, v30, v29
	v_add_u32_e32 v32, vcc, -1, v30
	v_subrev_u32_e32 v33, vcc, v27, v31
	v_cmp_ge_u32_e64 s[4:5], v31, v27
	v_cmp_ge_u32_e64 s[2:3], v33, v29
	v_add_u32_e32 v27, vcc, 1, v30
	s_and_b64 vcc, s[2:3], s[4:5]
	v_cndmask_b32_e32 v27, v30, v27, vcc
	v_cndmask_b32_e64 v27, v32, v27, s[4:5]
	v_xor_b32_e32 v29, v27, v28
	v_sub_u32_e32 v27, vcc, v29, v28
	v_cmp_gt_i32_e32 vcc, 1, v27
	s_and_b64 vcc, exec, vcc
	s_cbranch_vccnz BB2_348
; %bb.344:                              ; %.preheader7.i20.preheader
                                        ;   in Loop: Header=BB2_304 Depth=1
	v_not_b32_e32 v28, v28
	v_add_u32_e32 v28, vcc, v29, v28
	s_mov_b32 s2, 0
	v_mov_b32_e32 v29, v24
	s_branch BB2_346
BB2_345:                                ;   in Loop: Header=BB2_346 Depth=2
	v_mov_b32_e32 v29, 0
	ds_read_b32 v29, v29 offset:36
	s_mov_b64 s[4:5], 0
	s_andn2_b64 vcc, exec, s[4:5]
	s_add_i32 s2, s2, 1
	s_cbranch_vccz BB2_348
BB2_346:                                ; %.preheader7.i20
                                        ;   Parent Loop BB2_304 Depth=1
                                        ; =>  This Inner Loop Header: Depth=2
	s_waitcnt lgkmcnt(0)
	v_mul_lo_u32 v29, s2, v29
	v_add_u32_e32 v29, vcc, v29, v0
	v_add_u32_e32 v30, vcc, v29, v23
	v_add_u32_e32 v29, vcc, v29, v25
	v_lshlrev_b32_e32 v30, 2, v30
	v_lshlrev_b32_e32 v31, 2, v29
	ds_read2_b32 v[29:30], v30 offset0:16 offset1:17
	ds_read_b32 v32, v31 offset:64
	s_waitcnt lgkmcnt(1)
	v_add_u32_e32 v29, vcc, v30, v29
	v_lshrrev_b32_e32 v30, 31, v29
	v_add_u32_e32 v29, vcc, v29, v30
	v_ashrrev_i32_e32 v29, 1, v29
	s_waitcnt lgkmcnt(0)
	v_sub_u32_e32 v29, vcc, v32, v29
	v_cmp_eq_u32_e32 vcc, s2, v28
	s_and_b64 vcc, exec, vcc
	ds_write_b32 v31, v29 offset:64
	s_cbranch_vccz BB2_345
; %bb.347:                              ;   in Loop: Header=BB2_346 Depth=2
	s_mov_b64 s[4:5], -1
                                        ; implicit-def: $vgpr29
	s_andn2_b64 vcc, exec, s[4:5]
	s_add_i32 s2, s2, 1
	s_cbranch_vccnz BB2_346
BB2_348:                                ; %Flow499
                                        ;   in Loop: Header=BB2_304 Depth=1
	v_mul_lo_u32 v24, v27, v24
	v_subrev_u32_e32 v27, vcc, v24, v26
	v_subrev_u32_e32 v24, vcc, v27, v26
	v_ashrrev_i32_e32 v28, 31, v27
	v_cmp_lt_u64_e32 vcc, v[0:1], v[27:28]
	s_and_saveexec_b64 s[2:3], vcc
	s_cbranch_execz BB2_350
; %bb.349:                              ;   in Loop: Header=BB2_304 Depth=1
	v_add_u32_e32 v24, vcc, v24, v0
	v_add_u32_e32 v23, vcc, v24, v23
	v_add_u32_e32 v24, vcc, v24, v25
	v_lshlrev_b32_e32 v23, 2, v23
	v_lshlrev_b32_e32 v25, 2, v24
	ds_read2_b32 v[23:24], v23 offset0:16 offset1:17
	ds_read_b32 v26, v25 offset:64
	s_waitcnt lgkmcnt(1)
	v_add_u32_e32 v23, vcc, v24, v23
	v_lshrrev_b32_e32 v24, 31, v23
	v_add_u32_e32 v23, vcc, v23, v24
	v_ashrrev_i32_e32 v23, 1, v23
	s_waitcnt lgkmcnt(0)
	v_sub_u32_e32 v23, vcc, v26, v23
	ds_write_b32 v25, v23 offset:64
BB2_350:                                ;   in Loop: Header=BB2_304 Depth=1
	s_or_b64 exec, exec, s[2:3]
	v_mov_b32_e32 v25, 0
	s_waitcnt lgkmcnt(0)
	s_barrier
	s_waitcnt lgkmcnt(0)
	ds_read2_b32 v[23:24], v25 offset0:1 offset1:9
	s_waitcnt lgkmcnt(0)
	v_ashrrev_i32_e32 v28, 31, v24
	v_add_u32_e32 v26, vcc, v28, v24
	v_xor_b32_e32 v29, v26, v28
	v_cvt_f32_u32_e32 v26, v29
	v_rcp_iflag_f32_e32 v30, v26
	ds_read2_b32 v[26:27], v25 offset0:11 offset1:15
	v_mul_f32_e32 v25, s18, v30
	v_cvt_u32_f32_e32 v30, v25
	s_waitcnt lgkmcnt(0)
	v_mul_lo_u32 v23, v26, v23
	v_lshlrev_b32_e32 v26, 1, v26
	v_mul_hi_u32 v31, v30, v29
	v_add_u32_e32 v25, vcc, -1, v23
	v_mul_lo_u32 v23, v30, v29
	v_ashrrev_i32_e32 v32, 31, v25
	v_cmp_eq_u32_e64 s[2:3], 0, v31
	v_add_u32_e32 v31, vcc, v32, v25
	v_sub_u32_e32 v33, vcc, 0, v23
	v_cndmask_b32_e64 v23, v23, v33, s[2:3]
	v_mul_hi_u32 v23, v23, v30
	v_xor_b32_e32 v31, v31, v32
	v_xor_b32_e32 v28, v32, v28
	v_add_u32_e32 v33, vcc, v23, v30
	v_subrev_u32_e32 v23, vcc, v23, v30
	v_cndmask_b32_e64 v23, v23, v33, s[2:3]
	v_mul_hi_u32 v30, v23, v31
	v_add_u32_e32 v23, vcc, v27, v26
	v_mul_lo_u32 v27, v30, v29
	v_add_u32_e32 v32, vcc, -1, v30
	v_subrev_u32_e32 v33, vcc, v27, v31
	v_cmp_ge_u32_e64 s[4:5], v31, v27
	v_cmp_ge_u32_e64 s[2:3], v33, v29
	v_add_u32_e32 v27, vcc, 1, v30
	s_and_b64 vcc, s[2:3], s[4:5]
	v_cndmask_b32_e32 v27, v30, v27, vcc
	v_cndmask_b32_e64 v27, v32, v27, s[4:5]
	v_xor_b32_e32 v29, v27, v28
	v_sub_u32_e32 v27, vcc, v29, v28
	v_cmp_gt_i32_e32 vcc, 1, v27
	s_and_b64 vcc, exec, vcc
	s_cbranch_vccnz BB2_355
; %bb.351:                              ; %.preheader5.i22.preheader
                                        ;   in Loop: Header=BB2_304 Depth=1
	v_not_b32_e32 v28, v28
	v_add_u32_e32 v28, vcc, v29, v28
	s_mov_b32 s2, 0
	v_mov_b32_e32 v29, v24
	s_branch BB2_353
BB2_352:                                ;   in Loop: Header=BB2_353 Depth=2
	v_mov_b32_e32 v29, 0
	ds_read_b32 v29, v29 offset:36
	s_mov_b64 s[4:5], 0
	s_andn2_b64 vcc, exec, s[4:5]
	s_add_i32 s2, s2, 1
	s_cbranch_vccz BB2_355
BB2_353:                                ; %.preheader5.i22
                                        ;   Parent Loop BB2_304 Depth=1
                                        ; =>  This Inner Loop Header: Depth=2
	s_waitcnt lgkmcnt(0)
	v_mul_lo_u32 v29, s2, v29
	v_add_u32_e32 v29, vcc, v29, v0
	v_add_u32_e32 v30, vcc, v29, v23
	v_add_u32_e32 v29, vcc, v29, v26
	v_lshlrev_b32_e32 v30, 2, v30
	v_lshlrev_b32_e32 v31, 2, v29
	ds_read2_b32 v[29:30], v30 offset0:16 offset1:17
	ds_read_b32 v32, v31 offset:68
	s_waitcnt lgkmcnt(1)
	v_add_u32_e32 v29, vcc, v29, v30
	v_add_u32_e32 v29, vcc, 2, v29
	v_ashrrev_i32_e32 v30, 31, v29
	v_lshrrev_b32_e32 v30, 30, v30
	v_add_u32_e32 v29, vcc, v29, v30
	v_ashrrev_i32_e32 v29, 2, v29
	s_waitcnt lgkmcnt(0)
	v_add_u32_e32 v29, vcc, v29, v32
	v_cmp_eq_u32_e32 vcc, s2, v28
	s_and_b64 vcc, exec, vcc
	ds_write_b32 v31, v29 offset:68
	s_cbranch_vccz BB2_352
; %bb.354:                              ;   in Loop: Header=BB2_353 Depth=2
	s_mov_b64 s[4:5], -1
                                        ; implicit-def: $vgpr29
	s_andn2_b64 vcc, exec, s[4:5]
	s_add_i32 s2, s2, 1
	s_cbranch_vccnz BB2_353
BB2_355:                                ; %Flow496
                                        ;   in Loop: Header=BB2_304 Depth=1
	v_mul_lo_u32 v27, v27, v24
	v_or_b32_e32 v24, 1, v26
	v_subrev_u32_e32 v26, vcc, v27, v25
	v_subrev_u32_e32 v25, vcc, v26, v25
	v_ashrrev_i32_e32 v27, 31, v26
	v_cmp_lt_u64_e32 vcc, v[0:1], v[26:27]
	s_and_saveexec_b64 s[2:3], vcc
	s_cbranch_execz BB2_357
; %bb.356:                              ;   in Loop: Header=BB2_304 Depth=1
	v_add_u32_e32 v25, vcc, v25, v0
	v_add_u32_e32 v23, vcc, v25, v23
	v_add_u32_e32 v24, vcc, v25, v24
	v_lshlrev_b32_e32 v23, 2, v23
	v_lshlrev_b32_e32 v25, 2, v24
	ds_read2_b32 v[23:24], v23 offset0:16 offset1:17
	ds_read_b32 v26, v25 offset:64
	s_waitcnt lgkmcnt(1)
	v_add_u32_e32 v23, vcc, v23, v24
	v_add_u32_e32 v23, vcc, 2, v23
	v_ashrrev_i32_e32 v24, 31, v23
	v_lshrrev_b32_e32 v24, 30, v24
	v_add_u32_e32 v23, vcc, v23, v24
	v_ashrrev_i32_e32 v23, 2, v23
	s_waitcnt lgkmcnt(0)
	v_add_u32_e32 v23, vcc, v23, v26
	ds_write_b32 v25, v23 offset:64
BB2_357:                                ;   in Loop: Header=BB2_304 Depth=1
	s_or_b64 exec, exec, s[2:3]
	v_mov_b32_e32 v23, 0
	s_waitcnt lgkmcnt(0)
	s_barrier
	s_waitcnt lgkmcnt(0)
	ds_read_b32 v24, v23 offset:4
	s_waitcnt lgkmcnt(0)
	v_cmp_gt_i32_e32 vcc, 1, v24
	s_and_b64 vcc, exec, vcc
	s_cbranch_vccnz BB2_303
; %bb.358:                              ;   in Loop: Header=BB2_304 Depth=1
	v_add_u32_e32 v24, vcc, -1, v24
	v_lshrrev_b32_e32 v25, 1, v24
	v_add_u32_e32 v25, vcc, 1, v25
	v_cmp_gt_u32_e32 vcc, 6, v24
	ds_read_b32 v24, v23 offset:44
	s_and_b64 vcc, exec, vcc
	s_cbranch_vccnz BB2_377
; %bb.359:                              ;   in Loop: Header=BB2_304 Depth=1
	s_waitcnt lgkmcnt(0)
	v_mul_lo_u32 v23, v24, 9
	v_mul_lo_u32 v29, v24, 7
	v_mul_lo_u32 v30, v24, 6
	v_mul_lo_u32 v33, v24, 12
	v_add_u32_e32 v23, vcc, v20, v23
	v_lshlrev_b32_e32 v27, 2, v23
	v_lshlrev_b32_e32 v23, 3, v24
	v_add_u32_e32 v23, vcc, v20, v23
	v_lshlrev_b32_e32 v28, 2, v23
	v_add_u32_e32 v23, vcc, v20, v29
	v_lshlrev_b32_e32 v29, 2, v23
	v_mul_lo_u32 v23, v24, 5
	v_add_u32_e32 v30, vcc, v20, v30
	s_mov_b32 s20, 2
	v_lshlrev_b32_e32 v26, 5, v24
	v_add_u32_e32 v23, vcc, v20, v23
	v_lshlrev_b32_e32 v31, 2, v23
	v_lshlrev_b32_e32 v23, 2, v24
	v_add_u32_e32 v23, vcc, v20, v23
	v_lshlrev_b32_e32 v32, 2, v23
	v_lshlrev_b32_e32 v23, 1, v24
	v_add_u32_e32 v23, vcc, v20, v23
	v_lshlrev_b32_e32 v34, 2, v23
	v_and_b32_e32 v23, s15, v25
	v_lshlrev_b32_e32 v30, 2, v30
	v_add_u32_e32 v33, vcc, v21, v33
	v_sub_u32_e32 v35, vcc, 0, v23
	v_mov_b32_e32 v36, 0
	s_branch BB2_362
BB2_360:                                ; %Flow487
                                        ;   in Loop: Header=BB2_362 Depth=2
	s_or_b64 exec, exec, s[4:5]
BB2_361:                                ;   in Loop: Header=BB2_362 Depth=2
	s_or_b64 exec, exec, s[22:23]
	v_cndmask_b32_e64 v13, 0, 1, vcc
	v_add_u32_e32 v36, vcc, v36, v26
	v_add_u32_e32 v35, vcc, 4, v35
	v_cmp_eq_u32_e32 vcc, 0, v35
	s_and_b64 vcc, exec, vcc
	s_add_i32 s20, s20, 8
	s_cbranch_vccnz BB2_378
BB2_362:                                ;   Parent Loop BB2_304 Depth=1
                                        ; =>  This Inner Loop Header: Depth=2
	v_and_b32_e32 v13, s19, v13
	v_cmp_eq_u16_e64 s[2:3], 0, v13
	v_cmp_ne_u32_e64 s[4:5], 0, v22
	v_cmp_ne_u16_e32 vcc, 0, v13
	s_or_b64 s[2:3], s[2:3], s[4:5]
	v_mov_b32_e32 v37, 0
	s_and_saveexec_b64 s[4:5], s[2:3]
	s_cbranch_execz BB2_366
; %bb.363:                              ;   in Loop: Header=BB2_362 Depth=2
	v_add_u32_e64 v13, s[2:3], v36, v34
	v_ashrrev_i32_e32 v23, 31, v22
	ds_read_b32 v13, v13 offset:64
	v_lshlrev_b64 v[38:39], 2, v[22:23]
	v_mov_b32_e32 v23, s11
	v_add_u32_e64 v38, s[2:3], s10, v38
	v_addc_u32_e64 v39, s[2:3], v23, v39, s[2:3]
	s_waitcnt lgkmcnt(0)
	flat_store_dword v[38:39], v13
	v_add_u32_e64 v13, s[2:3], v22, v12
	v_cmp_ne_u32_e64 s[2:3], 0, v13
	s_xor_b64 s[22:23], vcc, -1
	s_or_b64 s[2:3], s[22:23], s[2:3]
	v_mov_b32_e32 v37, 0
	s_and_saveexec_b64 s[22:23], s[2:3]
	s_cbranch_execz BB2_365
; %bb.364:                              ;   in Loop: Header=BB2_362 Depth=2
	v_add_u32_e64 v22, s[2:3], v36, v33
	ds_read_b32 v37, v22 offset:64
	v_ashrrev_i32_e32 v23, 31, v12
	v_mov_b32_e32 v22, v12
	v_lshlrev_b64 v[22:23], 2, v[22:23]
	v_add_u32_e64 v22, s[2:3], v38, v22
	v_addc_u32_e64 v23, s[2:3], v39, v23, s[2:3]
	s_waitcnt lgkmcnt(0)
	flat_store_dword v[22:23], v37
	v_add_u32_e64 v37, s[2:3], v13, v14
BB2_365:                                ; %Flow490
                                        ;   in Loop: Header=BB2_362 Depth=2
	s_or_b64 exec, exec, s[22:23]
BB2_366:                                ;   in Loop: Header=BB2_362 Depth=2
	s_or_b64 exec, exec, s[4:5]
	v_cmp_ne_u32_e64 s[2:3], 0, v37
	s_xor_b64 s[4:5], vcc, -1
	s_or_b64 s[2:3], s[4:5], s[2:3]
	v_mov_b32_e32 v22, 0
	s_and_saveexec_b64 s[22:23], s[2:3]
	s_cbranch_execz BB2_370
; %bb.367:                              ;   in Loop: Header=BB2_362 Depth=2
	v_add_u32_e64 v13, s[2:3], v36, v32
	v_ashrrev_i32_e32 v38, 31, v37
	ds_read_b32 v13, v13 offset:64
	v_lshlrev_b64 v[22:23], 2, v[37:38]
	v_mov_b32_e32 v39, s11
	v_add_u32_e64 v38, s[2:3], s10, v22
	v_addc_u32_e64 v39, s[2:3], v39, v23, s[2:3]
	s_waitcnt lgkmcnt(0)
	flat_store_dword v[38:39], v13
	v_add_u32_e64 v13, s[2:3], v37, v12
	v_cmp_ne_u32_e64 s[2:3], 0, v13
	s_or_b64 s[2:3], s[4:5], s[2:3]
	v_mov_b32_e32 v22, 0
	s_and_saveexec_b64 s[24:25], s[2:3]
	s_cbranch_execz BB2_369
; %bb.368:                              ;   in Loop: Header=BB2_362 Depth=2
	v_add_u32_e64 v22, s[2:3], v36, v31
	ds_read_b32 v37, v22 offset:64
	v_ashrrev_i32_e32 v23, 31, v12
	v_mov_b32_e32 v22, v12
	v_lshlrev_b64 v[22:23], 2, v[22:23]
	v_add_u32_e64 v22, s[2:3], v38, v22
	v_addc_u32_e64 v23, s[2:3], v39, v23, s[2:3]
	s_waitcnt lgkmcnt(0)
	flat_store_dword v[22:23], v37
	v_add_u32_e64 v22, s[2:3], v13, v14
BB2_369:                                ; %Flow489
                                        ;   in Loop: Header=BB2_362 Depth=2
	s_or_b64 exec, exec, s[24:25]
BB2_370:                                ;   in Loop: Header=BB2_362 Depth=2
	s_or_b64 exec, exec, s[22:23]
	v_cmp_ne_u32_e64 s[2:3], 0, v22
	s_or_b64 s[2:3], s[4:5], s[2:3]
	v_mov_b32_e32 v37, 0
	s_and_saveexec_b64 s[22:23], s[2:3]
	s_cbranch_execz BB2_374
; %bb.371:                              ;   in Loop: Header=BB2_362 Depth=2
	v_add_u32_e64 v13, s[2:3], v36, v30
	v_ashrrev_i32_e32 v23, 31, v22
	ds_read_b32 v13, v13 offset:64
	v_lshlrev_b64 v[38:39], 2, v[22:23]
	v_mov_b32_e32 v23, s11
	v_add_u32_e64 v38, s[2:3], s10, v38
	v_addc_u32_e64 v39, s[2:3], v23, v39, s[2:3]
	s_waitcnt lgkmcnt(0)
	flat_store_dword v[38:39], v13
	v_add_u32_e64 v13, s[2:3], v22, v12
	v_cmp_ne_u32_e64 s[2:3], 0, v13
	s_or_b64 s[2:3], s[4:5], s[2:3]
	v_mov_b32_e32 v37, 0
	s_and_saveexec_b64 s[24:25], s[2:3]
	s_cbranch_execz BB2_373
; %bb.372:                              ;   in Loop: Header=BB2_362 Depth=2
	v_add_u32_e64 v22, s[2:3], v36, v29
	ds_read_b32 v37, v22 offset:64
	v_ashrrev_i32_e32 v23, 31, v12
	v_mov_b32_e32 v22, v12
	v_lshlrev_b64 v[22:23], 2, v[22:23]
	v_add_u32_e64 v22, s[2:3], v38, v22
	v_addc_u32_e64 v23, s[2:3], v39, v23, s[2:3]
	s_waitcnt lgkmcnt(0)
	flat_store_dword v[22:23], v37
	v_add_u32_e64 v37, s[2:3], v13, v14
BB2_373:                                ; %Flow488
                                        ;   in Loop: Header=BB2_362 Depth=2
	s_or_b64 exec, exec, s[24:25]
BB2_374:                                ;   in Loop: Header=BB2_362 Depth=2
	s_or_b64 exec, exec, s[22:23]
	v_cmp_ne_u32_e64 s[2:3], 0, v37
	s_or_b64 s[2:3], s[4:5], s[2:3]
	v_mov_b32_e32 v22, 0
	s_and_saveexec_b64 s[22:23], s[2:3]
	s_cbranch_execz BB2_361
; %bb.375:                              ;   in Loop: Header=BB2_362 Depth=2
	v_add_u32_e64 v13, s[2:3], v36, v28
	v_ashrrev_i32_e32 v38, 31, v37
	ds_read_b32 v13, v13 offset:64
	v_lshlrev_b64 v[22:23], 2, v[37:38]
	v_mov_b32_e32 v39, s11
	v_add_u32_e64 v38, s[2:3], s10, v22
	v_addc_u32_e64 v39, s[2:3], v39, v23, s[2:3]
	s_waitcnt lgkmcnt(0)
	flat_store_dword v[38:39], v13
	v_add_u32_e64 v13, s[2:3], v37, v12
	v_cmp_ne_u32_e64 s[2:3], 0, v13
	s_or_b64 s[2:3], s[4:5], s[2:3]
	v_mov_b32_e32 v22, 0
	s_and_saveexec_b64 s[4:5], s[2:3]
	s_cbranch_execz BB2_360
; %bb.376:                              ;   in Loop: Header=BB2_362 Depth=2
	v_add_u32_e64 v22, s[2:3], v36, v27
	ds_read_b32 v37, v22 offset:64
	v_ashrrev_i32_e32 v23, 31, v12
	v_mov_b32_e32 v22, v12
	v_lshlrev_b64 v[22:23], 2, v[22:23]
	v_add_u32_e64 v22, s[2:3], v38, v22
	v_addc_u32_e64 v23, s[2:3], v39, v23, s[2:3]
	s_waitcnt lgkmcnt(0)
	flat_store_dword v[22:23], v37
	v_add_u32_e64 v22, s[2:3], v13, v14
	s_branch BB2_360
BB2_377:                                ;   in Loop: Header=BB2_304 Depth=1
	s_mov_b32 s20, 2
BB2_378:                                ; %Flow492
                                        ;   in Loop: Header=BB2_304 Depth=1
	v_and_b32_e32 v25, 3, v25
	v_cmp_eq_u32_e32 vcc, 0, v25
	s_and_b64 vcc, exec, vcc
	s_cbranch_vccnz BB2_303
; %bb.379:                              ; %.preheader.i24.preheader
                                        ;   in Loop: Header=BB2_304 Depth=1
	s_waitcnt lgkmcnt(0)
	v_mul_lo_u32 v23, s20, v24
	v_lshlrev_b32_e32 v27, 3, v24
	v_add_u32_e32 v23, vcc, v20, v23
	v_lshlrev_b32_e32 v26, 2, v23
	s_branch BB2_382
BB2_380:                                ; %Flow484
                                        ;   in Loop: Header=BB2_382 Depth=2
	s_or_b64 exec, exec, s[22:23]
BB2_381:                                ;   in Loop: Header=BB2_382 Depth=2
	s_or_b64 exec, exec, s[4:5]
	v_cndmask_b32_e64 v13, 0, 1, vcc
	v_add_u32_e32 v25, vcc, -1, v25
	v_add_u32_e32 v26, vcc, v26, v27
	v_cmp_ne_u32_e32 vcc, 0, v25
	s_add_i32 s20, s20, 2
	s_and_b64 vcc, exec, vcc
	v_mov_b32_e32 v22, v23
	s_cbranch_vccz BB2_302
BB2_382:                                ; %.preheader.i24
                                        ;   Parent Loop BB2_304 Depth=1
                                        ; =>  This Inner Loop Header: Depth=2
	v_and_b32_e32 v13, s19, v13
	v_cmp_eq_u16_e64 s[2:3], 0, v13
	v_cmp_ne_u32_e64 s[4:5], 0, v22
	v_cmp_ne_u16_e32 vcc, 0, v13
	s_or_b64 s[2:3], s[2:3], s[4:5]
	v_mov_b32_e32 v23, 0
	s_and_saveexec_b64 s[4:5], s[2:3]
	s_cbranch_execz BB2_381
; %bb.383:                              ;   in Loop: Header=BB2_382 Depth=2
	v_ashrrev_i32_e32 v23, 31, v22
	v_lshlrev_b64 v[28:29], 2, v[22:23]
	ds_read_b32 v13, v26 offset:64
	v_mov_b32_e32 v23, s11
	v_add_u32_e64 v28, s[2:3], s10, v28
	v_addc_u32_e64 v29, s[2:3], v23, v29, s[2:3]
	v_add_u32_e64 v22, s[2:3], v22, v12
	v_cmp_ne_u32_e64 s[2:3], 0, v22
	s_xor_b64 s[22:23], vcc, -1
	s_or_b64 s[2:3], s[22:23], s[2:3]
	v_mov_b32_e32 v23, 0
	s_waitcnt lgkmcnt(0)
	flat_store_dword v[28:29], v13
	s_and_saveexec_b64 s[22:23], s[2:3]
	s_cbranch_execz BB2_380
; %bb.384:                              ;   in Loop: Header=BB2_382 Depth=2
	s_or_b32 s2, s20, 1
	v_mul_lo_u32 v23, s2, v24
	v_ashrrev_i32_e32 v13, 31, v12
	v_lshlrev_b64 v[30:31], 2, v[12:13]
	v_add_u32_e64 v13, s[2:3], v23, v15
	v_lshlrev_b32_e32 v13, 2, v13
	ds_read_b32 v13, v13 offset:64
	v_add_u32_e64 v28, s[2:3], v28, v30
	v_addc_u32_e64 v29, s[2:3], v29, v31, s[2:3]
	v_add_u32_e64 v23, s[2:3], v22, v14
	s_waitcnt lgkmcnt(0)
	flat_store_dword v[28:29], v13
	s_branch BB2_380
BB2_385:
	s_endpgm
.Lfunc_end2:
	.size	cl_fdwt53Kernel, .Lfunc_end2-cl_fdwt53Kernel
                                        ; -- End function
	.section	.AMDGPU.csdata
; Kernel info:
; codeLenInByte = 24180
; NumSgprs: 40
; NumVgprs: 42
; ScratchSize: 0
; MemoryBound: 0
; FloatMode: 192
; IeeeMode: 1
; LDSByteSize: 8796 bytes/workgroup (compile time only)
; SGPRBlocks: 4
; VGPRBlocks: 10
; NumSGPRsForWavesPerEU: 40
; NumVGPRsForWavesPerEU: 42
; Occupancy: 5
; WaveLimiterHint : 1
; COMPUTE_PGM_RSRC2:USER_SGPR: 6
; COMPUTE_PGM_RSRC2:TRAP_HANDLER: 0
; COMPUTE_PGM_RSRC2:TGID_X_EN: 1
; COMPUTE_PGM_RSRC2:TGID_Y_EN: 1
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
  - Name:            c_CopySrcToComponents
    SymbolName:      'c_CopySrcToComponents@kd'
    Language:        OpenCL C
    LanguageVersion: [ 1, 2 ]
    Args:
      - Name:            d_r
        TypeName:        'int*'
        Size:            8
        Align:           8
        ValueKind:       GlobalBuffer
        ValueType:       I32
        AddrSpaceQual:   Global
        AccQual:         Default
      - Name:            d_g
        TypeName:        'int*'
        Size:            8
        Align:           8
        ValueKind:       GlobalBuffer
        ValueType:       I32
        AddrSpaceQual:   Global
        AccQual:         Default
      - Name:            d_b
        TypeName:        'int*'
        Size:            8
        Align:           8
        ValueKind:       GlobalBuffer
        ValueType:       I32
        AddrSpaceQual:   Global
        AccQual:         Default
      - Name:            cl_d_src
        TypeName:        'uchar*'
        Size:            8
        Align:           8
        ValueKind:       GlobalBuffer
        ValueType:       U8
        AddrSpaceQual:   Global
        AccQual:         Default
      - Name:            pixels
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
      GroupSegmentFixedSize: 768
      PrivateSegmentFixedSize: 0
      KernargSegmentAlign: 8
      WavefrontSize:   64
      NumSGPRs:        22
      NumVGPRs:        7
      MaxFlatWorkGroupSize: 256
  - Name:            c_CopySrcToComponent
    SymbolName:      'c_CopySrcToComponent@kd'
    Language:        OpenCL C
    LanguageVersion: [ 1, 2 ]
    Args:
      - Name:            d_c
        TypeName:        'int*'
        Size:            8
        Align:           8
        ValueKind:       GlobalBuffer
        ValueType:       I32
        AddrSpaceQual:   Global
        AccQual:         Default
      - Name:            cl_d_src
        TypeName:        'uchar*'
        Size:            8
        Align:           8
        ValueKind:       GlobalBuffer
        ValueType:       U8
        AddrSpaceQual:   Global
        AccQual:         Default
      - Name:            pixels
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
      GroupSegmentFixedSize: 256
      PrivateSegmentFixedSize: 0
      KernargSegmentAlign: 8
      WavefrontSize:   64
      NumSGPRs:        11
      NumVGPRs:        5
      MaxFlatWorkGroupSize: 256
  - Name:            cl_fdwt53Kernel
    SymbolName:      'cl_fdwt53Kernel@kd'
    Language:        OpenCL C
    LanguageVersion: [ 1, 2 ]
    Args:
      - Name:            in
        TypeName:        'int*'
        Size:            8
        Align:           8
        ValueKind:       GlobalBuffer
        ValueType:       I32
        AddrSpaceQual:   Global
        AccQual:         Default
        IsConst:         true
      - Name:            out
        TypeName:        'int*'
        Size:            8
        Align:           8
        ValueKind:       GlobalBuffer
        ValueType:       I32
        AddrSpaceQual:   Global
        AccQual:         Default
      - Name:            sx
        TypeName:        int
        Size:            4
        Align:           4
        ValueKind:       ByValue
        ValueType:       I32
        AccQual:         Default
      - Name:            sy
        TypeName:        int
        Size:            4
        Align:           4
        ValueKind:       ByValue
        ValueType:       I32
        AccQual:         Default
      - Name:            steps
        TypeName:        int
        Size:            4
        Align:           4
        ValueKind:       ByValue
        ValueType:       I32
        AccQual:         Default
      - Name:            WIN_SIZE_X
        TypeName:        int
        Size:            4
        Align:           4
        ValueKind:       ByValue
        ValueType:       I32
        AccQual:         Default
      - Name:            WIN_SIZE_Y
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
      GroupSegmentFixedSize: 8796
      PrivateSegmentFixedSize: 0
      KernargSegmentAlign: 8
      WavefrontSize:   64
      NumSGPRs:        40
      NumVGPRs:        42
      MaxFlatWorkGroupSize: 256
...

	.end_amd_amdgpu_hsa_metadata
